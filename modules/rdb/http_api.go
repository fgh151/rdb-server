package rdb

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
	admin.HandleFunc("/rdb", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/rdb", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/rdb/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/rdb/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/rdb/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandler
}

// list godoc
// @Summary      List rdb
// @Description  List rdb
// @Tags         RDB
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   Rdb
//
// @Router       /admin/rdb [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Rdb{}, []string{}, r, w)
}

// item godoc
// @Summary      Rdbs item
// @Description  Rdb detail info
// @Tags         RDB
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Rdb id" id
// @Security bearerAuth
// @Success      200  {object}   Rdb
//
// @Router       /admin/rdb/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Rdb{}, w, r)
}

// update
// @Summary      Update rdb
// @Description  Update rdb
// @Tags         RDB
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     Rdb  true  "Rdb info" true
// @Param        id    path     string  true  "Rdb id" true
// @Success      200 {object} Rdb */
// @Security bearerAuth
//
// @Router       /admin/rdb/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	m, err := Rdb{}.GetById(vars["id"])

	if err != nil {
		w.WriteHeader(404)
		return
	}

	t := m.(Rdb)

	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	db.MetaDb.GetConnection().Save(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// deleteItem godoc
// @Summary      Delete rdb
// @Description  Delete rdb
// @Tags         RDB
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Rdb id" string
// @Success      204
//
// @Router       /admin/rdb/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Rdb{}, w, r)
}

// create
// @Summary      Create rdb
// @Description  Create rdb
// @Tags         RDB
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        rdb    body     Rdb  true  "rdb info" true
// @Success      200 {object} Rdb */
// @Security bearerAuth
//
// @Router       /admin/rdb [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t Rdb
	t.Id, _ = uuid.NewUUID()

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	db.MetaDb.GetConnection().Create(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
