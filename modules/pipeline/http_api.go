package pipeline

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
	admin.HandleFunc("/pl", ListPipeline).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/pl", CreatePipeline).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", PipelineItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", DeletePipeline).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", UpdatePipeline).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler
}

// ListPipeline godoc
// @Summary      List pipelines
// @Description  List pipelines
// @Tags         Pipeline
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   pipeline.Pipeline
//
// @Router       /admin/pl [get]
func ListPipeline(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Pipeline{}, []string{"id"}, r, w)
}

// CreatePipeline
// @Summary      Create pipeline
// @Description  Create pipeline
// @Tags         Pipeline
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        pl    body     pipeline.Pipeline  true  "Pipeline info" true
// @Success      200 {object} pipeline.Pipeline
// @Security bearerAuth
//
// @Router       /admin/pl [post]
func CreatePipeline(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := Pipeline{}

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

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// PipelineItem godoc
// @Summary      Pipeline info
// @Description  Pipeline detail info
// @Tags         Pipeline
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Pipeline id" id
// @Security bearerAuth
// @Success      200  {object}   pipeline.Pipeline
//
// @Router       /admin/pl/{id} [get]
func PipelineItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Pipeline{}, w, r)
}

// DeletePipeline godoc
// @Summary      Delete pipeline
// @Description  Delete pipeline
// @Tags         Pipeline
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Pipeline id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/pl/{id} [delete]
func DeletePipeline(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Pipeline{}, w, r)
}

// UpdatePipeline
// @Summary      Update pipeline
// @Description  Update pipeline
// @Tags         Pipeline
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     pipeline.Pipeline  true  "Pipeline info" true
// @Param        id    path     string  true  "Pipeline id" id
// @Success      200 {object} pipeline.Pipeline
// @Security bearerAuth
//
// @Router       /admin/pl/{id} [put]
func UpdatePipeline(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = Pipeline{}.GetById(vars["id"]).(Pipeline)
	newm := Pipeline{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
