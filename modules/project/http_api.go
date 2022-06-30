package project

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
	admin.HandleFunc("/projects", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/projects", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/projects/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/projects/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/projects/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandler
}

// list godoc
// @Summary      List projects
// @Description  List projects
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   Project
//
// @Router       /admin/projects [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Project{}, []string{}, r, w)
}

// item godoc
// @Summary      Projects item
// @Description  Project detail info
// @Tags         Projects
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Project id" id
// @Security bearerAuth
// @Success      200  {object}   Project
//
// @Router       /admin/projects/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Project{}, w, r)
}

// update
// @Summary      Update project
// @Description  Update project
// @Tags         Projects
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     Project  true  "Project info" true
// @Param        id    path     string  true  "Project id" true
// @Success      200 {object} Project */
// @Security bearerAuth
//
// @Router       /admin/projects/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	m, err := Project{}.GetById(vars["id"])

	if err != nil {
		w.WriteHeader(404)
		return
	}

	t := m.(Project)

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
// @Summary      Delete project
// @Description  Delete project
// @Tags         Projects
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Project id" string
// @Success      204
//
// @Router       /admin/projects/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Project{}, w, r)
}

// create
// @Summary      Create project
// @Description  Create project
// @Tags         Projects
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        project    body     Project  true  "project info" true
// @Success      200 {object} Project */
// @Security bearerAuth
//
// @Router       /admin/projects [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t Project
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
