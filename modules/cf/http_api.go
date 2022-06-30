package cf

import (
	err2 "db-server/err"
	"db-server/server"
	"db-server/server/db"
	"db-server/utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/cf", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/cf", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/cf/{id}/log", logs).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandler
}

func AddApiRoutes(api *mux.Router) {

	api.HandleFunc("/cf/{id}/run", CfRun).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	api.HandleFunc("/cf/{id}/run/{rid}", CfRunLog).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
}

// list godoc
// @Summary      List cloud functions
// @Description  List cloud functions
// @Tags         Cloud functions
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   cf.CloudFunction
//
// @Router       /admin/cf [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(CloudFunction{}, []string{"id"}, r, w)
}

// create
// @Summary      Create cloud function
// @Description  Create cloud function
// @Tags         Cloud functions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        cf    body     cf.CloudFunction  true  "Function info" true
// @Success      200 {object} cf.CloudFunction
// @Security bearerAuth
//
// @Router       /admin/cf [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := CloudFunction{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	id, err := uuid.NewUUID()
	model.Id = id

	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}

	file, _, err := r.FormFile("dockerarc")
	if err == nil {
		uri, err := GetContainerUri(model.Container)
		err2.DebugErr(err)

		go func() {
			err := server.BuildDockerImage(file, []string{uri.Vendor + "/" + uri.Image})
			err2.DebugErr(err)
		}()
	}

	db.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// item godoc
// @Summary      Function info
// @Description  Function detail info
// @Tags         Cloud functions
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "cf id" id
// @Security bearerAuth
// @Success      200  {object}   cf.CloudFunction
//
// @Router       /admin/cf/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(CloudFunction{}, w, r)
}

// logs godoc
// @Summary      Logs
// @Description  Cloud function logs
// @Tags         Cloud functions
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Fuc id" id
// @Security bearerAuth
// @Success      200  {object}   cf.CloudFunctionLog
//
// @Router       /admin/cf/{id}/log [get]
func logs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	f := CloudFunction{}.GetById(vars["id"]).(CloudFunction)

	l, o, s, or := utils.GetPagination(r)
	arr := ListCfLog(f.Id, l, o, s, or)
	total := LogsTotal(f.Id)
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, _ := json.Marshal(arr)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// deleteItem godoc
// @Summary      Delete cloud function
// @Description  Delete cloud function
// @Tags         Cloud functions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "cf id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/cf/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(CloudFunction{}, w, r)
}

// update
// @Summary      Update cloud function
// @Description  Update cloud function
// @Tags         Cloud functions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     ds.DataSource  true  "Source info" true
// @Param        id    path     string  true  "Function id" id
// @Success      200 {object} ds.DataSource
// @Security bearerAuth
//
// @Router       /admin/cf/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)

	var projectId, _ = uuid.Parse(r.FormValue("project_id"))

	uri, err := GetContainerUri(r.FormValue("container"))
	file, _, err := r.FormFile("dockerarc")
	if err == nil {
		err2.DebugErr(err)

		go func() {
			err := server.BuildDockerImage(file, []string{uri.Vendor + "/" + uri.Image})
			err2.DebugErr(err)
		}()
	} else {
		log.Debug(err)
		server.PullDockerImage(uri.Vendor + "/" + uri.Image)
	}

	db.MetaDb.GetConnection().Table("cloud_functions").Where("id = ?", vars["id"]).Updates(
		map[string]interface{}{
			"title":      r.FormValue("title"),
			"project_id": projectId,
			"container":  r.FormValue("container"),
			"params":     r.FormValue("params"),
			"env":        r.FormValue("env"),
		},
	)

	resp, _ := json.Marshal(CloudFunction{}.GetById(vars["id"]))
	w.WriteHeader(200)
	_, err = w.Write(resp)
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
	cfu := CloudFunction{}.GetById(vars["id"]).(CloudFunction)

	id, _ := uuid.NewUUID()

	go cfu.Run(id)
	m := make(map[string]string)
	m["id"] = id.String()

	resp, _ := json.Marshal(m)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
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
// @Success      200 {object} cf.CloudFunctionLog
//
// @Router       /api/cf/{id}/run/{rid} [get]
func CfRunLog(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var logModel CloudFunctionLog

	conn := db.MetaDb.GetConnection()

	conn.First(&logModel, "id = ? AND function_id = ?", vars["rid"], vars["id"])

	resp, _ := json.Marshal(logModel)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}
