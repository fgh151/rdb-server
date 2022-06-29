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
	admin.HandleFunc("/projects", ListProjects).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/projects", CreateProject).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/projects/{id}", ProjectItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/projects/{id}", DeleteProject).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/projects/{id}", UpdateProject).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler
}

// ListProjects godoc
// @Summary      List projects
// @Description  List projects
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   Project
//
// @Router       /admin/projects [get]
func ListProjects(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Project{}, []string{}, r, w)
}

// ProjectItem godoc
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
func ProjectItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Project{}, w, r)
}

// UpdateProject
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
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = Project{}.GetById(vars["id"]).(Project)

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

// DeleteProject godoc
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
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Project{}, w, r)
}

// CreateProject
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
func CreateProject(w http.ResponseWriter, r *http.Request) {
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
