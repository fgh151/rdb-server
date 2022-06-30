package cron

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
)

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/cron", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/cron", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandlerx
}

// list godoc
// @Summary      List cron jobs
// @Description  List cron jobs
// @Tags         Cron
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   cron.CronJob
//
// @Router       /admin/cron [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(CronJob{}, []string{}, r, w)
}

// create
// @Summary      Create cron job
// @Description  Create cron job
// @Tags         Cron
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        cron    body     cron.CronJob  true  "Push info" true
// @Success      200 {object} cron.CronJob
// @Security bearerAuth
//
// @Router       /admin/cron [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := CronJob{}

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
	db.MetaDb.GetConnection().Create(&model)

	model.Schedule(server.Cron.GetScheduler())

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// item godoc
// @Summary      Cron job info
// @Description  Cron job detail info
// @Tags         Cron
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Push id" id
// @Security bearerAuth
// @Success      200  {object}   cron.CronJob
//
// @Router       /admin/cron/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(CronJob{}, w, r)
}

// deleteItem godoc
// @Summary      Delete cron job
// @Description  Delete cron job
// @Tags         Cron
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Cron id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/cron/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(CronJob{}, w, r)
}

// update
// @Summary      Update cron job
// @Description  Update cron job
// @Tags         Cron
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     cron.CronJob  true  "Cron job info" true
// @Param        id    path     string  true  "Cron id"
// @Success      200 {object} cron.CronJob
// @Security bearerAuth
//
// @Router       /admin/cron/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	m, err := CronJob{}.GetById(vars["id"])

	if err != nil {
		w.WriteHeader(404)
		return
	}

	exist := m.(CronJob)

	newm := CronJob{}

	err = json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt
	db.MetaDb.GetConnection().Save(&newm)

	c := server.Cron
	c.GetScheduler().Remove(exist.CronId)
	newm.Schedule(c.GetScheduler())

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
