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
	admin.HandleFunc("/rdb", ListRdbs).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/rdb", CreateRdb).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/rdb/{id}", RdbItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/rdb/{id}", DeleteRdb).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/rdb/{id}", UpdateRdb).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler
}

// ListRdbs godoc
// @Summary      List rdb
// @Description  List rdb
// @Tags         RDB
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   Rdb
//
// @Router       /admin/rdb [get]
func ListRdbs(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Rdb{}, []string{}, r, w)
}

// RdbItem godoc
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
func RdbItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Rdb{}, w, r)
}

// UpdateRdb
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
func UpdateRdb(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = Rdb{}.GetById(vars["id"]).(Rdb)

	err := json.NewDecoder(r.Body).Decode(&t)
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

// DeleteRdb godoc
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
func DeleteRdb(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Rdb{}, w, r)
}

// CreateRdb
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
func CreateRdb(w http.ResponseWriter, r *http.Request) {
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
