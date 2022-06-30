package user

import (
	err2 "db-server/err"
	"db-server/server/db"
	"db-server/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/users", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/users", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/users/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/users/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/users/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandler
}

func AddPublicApiRoutes(r *mux.Router) {
	r.HandleFunc("/admin/auth", adminAuth).Methods(http.MethodPost, http.MethodOptions)                   // each request calls PushHandler
	r.HandleFunc("/api/user/auth", apiAuth).Methods(http.MethodPost, http.MethodOptions)                  // each request calls PushHandler
	r.HandleFunc("/api/user/register", register).Methods(http.MethodPost, http.MethodOptions)             // each request calls PushHandler
	r.HandleFunc("/api/device/register", pushDeviceRegister).Methods(http.MethodPost, http.MethodOptions) // each request calls PushHandler
}

func AddApiRoutes(api *mux.Router) {
	api.HandleFunc("/user/me", ApiMe).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
}

// list godoc
// @Summary      List users
// @Description  List users
// @Tags         User
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   User
//
// @Router       /admin/users [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(User{}, []string{"id", "email", "admin", "active"}, r, w)
}

// create
// @Summary      Create user
// @Description  Create user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        user    body     CreateUserForm  true  "User info" true
// @Success      200 {object} User
// @Security bearerAuth
//
// @Router       /admin/users [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t CreateUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	usr := t.Save()
	resp, _ := json.Marshal(usr)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// item godoc
// @Summary      User info
// @Description  User detail info
// @Tags         User
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "User id" gg
// @Security bearerAuth
// @Success      200  {object}   User
//
// @Router       /admin/users/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(User{}, w, r)
}

// deleteItem godoc
// @Summary      Delete user
// @Description  Delete user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "User id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/users/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(User{}, w, r)
}

// update
// @Summary      Update user
// @Description  Update user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     User  true  "User info" true
// @Param        id    path     string  true  "User info" id
// @Success      200 {object} User
// @Security bearerAuth
//
// @Router       /admin/users/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = User{}.GetById(vars["id"]).(User)
	newm := User{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt
	newm.LastLogin = exist.LastLogin
	newm.PasswordHash = exist.PasswordHash

	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// adminAuth godoc
// @Summary      Login to admin
// @Description  Authenticate in admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   User
//
// @Router       /admin/auth [post]
func adminAuth(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var l LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	usr, err := l.AdminLogin()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(usr)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

	return
}

// apiAuth godoc
// @Summary      Login via api
// @Description  Authenticate via api
// @Tags         User
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   User
//
// @Router       /api/user/auth [post]
func apiAuth(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var l LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	usr, err := l.ApiLogin()

	if err != nil {
		err2.DebugErr(err)
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(usr)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

	return
}

// register godoc
// @Summary      Register via api
// @Description  Register via api
// @Tags         User
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   User
//
// @Router       /api/user/register [post]
func register(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t CreateUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	usr := t.Save()
	now := time.Now()
	usr.LastLogin = &now
	db.MetaDb.GetConnection().Save(&usr)

	resp, _ := json.Marshal(usr)
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
// @Success      200  {object}   User
//
// @Router       /api/user/me [get]
func ApiMe(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	usr, err := GetUserFromRequest(r)
	resp, _ := json.Marshal(usr)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// pushDeviceRegister
// @Summary      Register device
// @Description  Register device to receive push
// @Tags         Push messages
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        device    body     user.UserDevice  true  "Device info" true
// @Success      200 {object} user.UserDevice
//
// @Router       /api/device/register [post]
func pushDeviceRegister(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var result map[string]string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)

	err2.DebugErr(err)

	var device UserDevice
	conn := db.MetaDb.GetConnection()
	res := conn.First(&device, "device_token = ?", result["device_token"])

	log.Debug(res.RowsAffected)

	if res.RowsAffected == 0 {

		id, err := uuid.NewUUID()
		err2.DebugErr(err)
		userId, err := uuid.Parse(result["user_id"])
		err2.DebugErr(err)

		device = UserDevice{
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

// GetUserFromRequest Fetch user model from request
func GetUserFromRequest(r *http.Request) (User, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	var usr = User{}

	if len(splitToken) < 2 {
		return usr, nil
	}

	reqToken = splitToken[1]

	tx := db.MetaDb.GetConnection().Table("user").Find(&usr, "token = ? ", reqToken)

	if tx.RowsAffected < 1 {
		return usr, errors.New("invalid credentials")
	}

	return usr, nil
}
