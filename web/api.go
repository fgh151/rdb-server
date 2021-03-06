package web

import (
	"db-server/auth"
	err2 "db-server/err"
	"db-server/models"
	"db-server/oauth"
	"db-server/server"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"time"
)

// ApiAuth godoc
// @Summary      Login via api
// @Description  Authenticate via api
// @Tags         User
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   models.User
//
// @Router       /api/user/auth [post]
func ApiAuth(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var l models.LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := l.ApiLogin()

	if err != nil {
		err2.DebugErr(err)
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

	return
}

// ApiOAuthLink godoc
// @Summary      OAuth link
// @Description  Get link for oauth
// @Tags         OAuth
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        provider    path     string  true  "Provider name"
// @Param        db-key    header     string  true  "Auth key" gg
// @Success      200  {string} string
//
// @Router       /api/user/oauth/{provider}/link [get]
func ApiOAuthLink(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	provider := vars["provider"]

	rKey := r.Header.Get("db-key")
	p, err := models.Project{}.GetByKey(rKey)

	if err != nil {
		payload := map[string]string{"code": "not acceptable", "message": err.Error()}
		sendResponse(w, 500, payload, nil)
	}

	client, _ := oauth.GetClient(provider, p.Id)

	url := client.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	w.WriteHeader(200)
	w.Write([]byte(url))
}

// ApiOAuthCode godoc
// @Summary      OAuth user
// @Description  Get user by oauth code
// @Tags         OAuth
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        provider    path     string  true  "Provider name"
// @Param        code    path     string  true  "Code"
// @Param        db-key    header     string  true  "Auth key" gg
// @Success      200  {object} models.User
//
// @Router       /api/user/oauth/{provider}/{code} [get]
func ApiOAuthCode(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	provider := vars["provider"]
	code := vars["code"]

	rKey := r.Header.Get("db-key")
	p, err := models.Project{}.GetByKey(rKey)

	if err != nil {
		payload := map[string]string{"code": "not acceptable", "message": err.Error()}
		sendResponse(w, 500, payload, nil)
	}

	client, _ := oauth.GetClient(provider, p.Id)

	u, err := client.GetUserByCode(code)

	if err != nil {
		log.Debug(err)
		payload := map[string]string{"code": "not acceptable", "message": err.Error()}
		sendResponse(w, 500, payload, nil)
		return
	}

	user := u.GetUser()
	user.UpdateLastLogin()

	rresp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err = w.Write(rresp)
	err2.DebugErr(err)

	return
}

// ApiRegister godoc
// @Summary      Register via api
// @Description  Register via api
// @Tags         User
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   models.User
//
// @Router       /api/user/register [post]
func ApiRegister(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t models.CreateUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user := t.Save()
	now := time.Now()
	user.LastLogin = &now
	server.MetaDb.GetConnection().Save(&user)

	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// ApiMe godoc
// @Summary      User info
// @Description  Get current user info
// @Tags         User
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   models.User
//
// @Router       /api/user/me [get]
func ApiMe(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	user := auth.GetUserFromRequest(r)
	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// ApiConfigItem godoc
// @Summary      Config
// @Description  Get config by id
// @Tags         Config manager
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Config id" id
// @Param        db-key    header     string  true  "Auth key" gg
// @Success      200  {array}   interface{}
//
// @Router       /config/{id} [get]
func ApiConfigItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	model := models.Config{}.GetById(vars["id"]).(models.Config)

	rKey := r.Header.Get("db-key")

	if !validateKey(model.Project.Key, rKey) {
		Send403Error(w, "db-key not Valid")
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, err := w.Write([]byte(model.Body))
		err2.DebugErr(err)
	}
}

// DSEItem godoc
// @Summary      Get item
// @Description  Get data source by id
// @Tags         Data source
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        db-key    header     string  true  "Auth key" gg
// @Param        id    path     string  true  "Source id"
// @Success      200  {object}   models.Project
//
// @Router       /dse/{id} [get]
func DSEItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	model := models.DataSourceEndpoint{}.GetById(vars["id"]).(models.DataSourceEndpoint)

	arr := model.List(10, 0, "id", "ASC", make(map[string]interface{}))
	total := model.Total()
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, _ := json.Marshal(arr)

	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// CfRun godoc
// @Summary      Run function
// @Description  Run function with id
// @Tags         Cloud functions
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        db-key    header     string  false  "Auth key" gg
// @Param        id    path     string  true  "Function id" gg
// @Success      200
//
// @Router       /api/cf/{id}/run [get]
func CfRun(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	cf := models.CloudFunction{}.GetById(vars["id"]).(models.CloudFunction)

	id, _ := uuid.NewUUID()

	go cf.Run(id)
	m := make(map[string]string)
	m["id"] = id.String()

	resp, _ := json.Marshal(m)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// PushRun
// @Summary      Send push
// @Description  Send push with id
// @Tags         Push messages
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        db-key    header     string  true  "Auth key" true
// @Param        id    path     string  true  "Push id"
// @Success      200
//
// @Router       /api/push/{id}/run [get]
func PushRun(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	message := models.PushMessage{}.GetById(vars["id"]).(models.PushMessage)
	go message.Send()
	w.WriteHeader(200)
}

// CfRunLog
// @Summary      List logs
// @Description  List logs of function run
// @Tags         Cloud functions
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        db-key    header     string  false  "Auth key" true
// @Param        id    path     string  true  "Function id"
// @Param        rid    path     string  true  "Run id"
// @Success      200 {object} models.CloudFunctionLog
//
// @Router       /api/cf/{id}/run/{rid} [get]
func CfRunLog(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var logModel models.CloudFunctionLog

	conn := server.MetaDb.GetConnection()

	conn.First(&logModel, "id = ? AND function_id = ?", vars["rid"], vars["id"])

	resp, _ := json.Marshal(logModel)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// PushDeviceRegister
// @Summary      Register device
// @Description  Register device to receive push
// @Tags         Push messages
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        device    body     models.UserDevice  true  "Device info" true
// @Success      200 {object} models.UserDevice
//
// @Router       /api/device/register [post]
func PushDeviceRegister(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var result map[string]string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)

	err2.DebugErr(err)

	var device models.UserDevice
	conn := server.MetaDb.GetConnection()
	res := conn.First(&device, "device_token = ?", result["device_token"])

	log.Debug(res.RowsAffected)

	if res.RowsAffected == 0 {

		id, err := uuid.NewUUID()
		err2.DebugErr(err)
		userId, err := uuid.Parse(result["user_id"])
		err2.DebugErr(err)

		device = models.UserDevice{
			Id:          id,
			DeviceToken: result["device_token"],
			UserId:      userId,
			Device:      result["device"],
		}

		log.Debug("Create user device " + device.Id.String())
		conn.Create(&device)

	} else {
		log.Debug("Update user device " + device.Id.String())
		device.Device = result["device"]
		device.UserId, _ = uuid.Parse(result["user_id"])
		device.DeviceToken = result["device_token"]
		conn.Save(&device)
	}

	resp, _ := json.Marshal(device)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
