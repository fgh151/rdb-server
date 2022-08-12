package plugin

import (
	err2 "db-server/err"
	"db-server/server/db"
	"db-server/utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/plugin", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/plugin", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/plugin/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/plugin/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/plugin/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandlerx
}

// list godoc
// @Summary      List plugins
// @Description  List plugins
// @Tags         Plugin
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   plugin.Plugin
//
// @Router       /admin/plugin [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Plugin{}, []string{}, r, w)
}

// create
// @Summary      Create plugin
// @Description  Create plugin
// @Tags         Plugin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        plugin    body     plugin.Plugin  true  "Plugin info" true
// @Success      200 {object} plugin.Plugin
// @Security bearerAuth
//
// @Router       /admin/plugin [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := Plugin{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	model.Id, err = uuid.NewUUID()
	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	db.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// item godoc
// @Summary      Plugin info
// @Description  Plugin info
// @Tags         Plugin
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Plugin id" id
// @Security bearerAuth
// @Success      200  {object}   plugin.Plugin
//
// @Router       /admin/plugin/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Plugin{}, w, r)
}

// deleteItem godoc
// @Summary      Delete plugin
// @Description  Delete plugin
// @Tags         Plugin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "plugin id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/plugin/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Plugin{}, w, r)
}

// update
// @Summary      Update plugin
// @Description  Update plugin
// @Tags         Plugin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     plugin.Plugin  true  "plugin info" true
// @Param        id    path     string  true  "plugin id"
// @Success      200 {object} plugin.Plugin
// @Security bearerAuth
//
// @Router       /admin/plugin/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	newm := Plugin{}

	err := json.NewDecoder(r.Body).Decode(&newm)
	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
