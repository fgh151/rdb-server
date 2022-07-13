package ds

import (
	err2 "db-server/err"
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
	admin.HandleFunc("/ds/dse/{dsId}", listDse).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}", createDse).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", dseItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", deleteDse).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", updateDse).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/ds", listDs).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds", createDs).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", dsItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", deleteDs).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", updateDs).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler
}

func AddPublicApiRoutes(r *mux.Router) {
	r.HandleFunc("/dse/{id}", publicDesItem).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
}

// listDs godoc
// @Summary      List data sources
// @Description  List data sources
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   DataSource
//
// @Router       /admin/ds [get]
func listDs(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(DataSource{}, []string{}, r, w)
}

// listDse godoc
// @Summary      List data source endpoints
// @Description  List data source endpoints
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Param        dsIid path    string  true  "Data source id" id
// @Success      200  {array}   DataSourceEndpoint
//
// @Router       /admin/ds/{dsIid}/dse [get]
func listDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	l, o, or, so := utils.GetPagination(r)
	f := utils.FormatQuery(r, []string{"data_source_id"})

	arr, _ := DataSourceEndpoint{}.List(l, o, so, or, f)
	total := len(arr)
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", strconv.Itoa(total))

	resp, _ := json.Marshal(arr)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// createDs
// @Summary      Create data source
// @Description  Create data source
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        ds    body     DataSource  true  "Data source info" true
// @Success      200 {object} DataSource
// @Security bearerAuth
//
// @Router       /admin/ds [post]
func createDs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := DataSource{}

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

// createDse
// @Summary      Create data source endpoint
// @Description  Create data source endpoint
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        dse    body     DataSourceEndpoint  true  "Data source info" true
// @Param        dsId    path     string  true  "Data source id" id
// @Success      200 {object} DataSourceEndpoint
// @Security bearerAuth
//
// @Router       /admin/ds/dse/{dsId} [post]
func createDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	dsUuid, err := uuid.Parse(vars["dsId"])
	err2.DebugErr(err)

	if err != nil {
		payload := map[string]string{"code": "not acceptable", "message": "Wrong data source id"}
		w.WriteHeader(500)
		resp, _ := json.Marshal(payload)
		_, err = w.Write(resp)
		return
	}

	model := DataSourceEndpoint{
		DataSourceId: dsUuid,
	}

	err = json.NewDecoder(r.Body).Decode(&model)
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

// dsItem godoc
// @Summary      Data source info
// @Description  Data source detail info
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Source id" gg
// @Security bearerAuth
// @Success      200  {object}   DataSource
//
// @Router       /admin/ds/{id} [get]
func dsItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(DataSource{}, w, r)
}

// dseItem godoc
// @Summary      Data source endpoint info
// @Description  Data source endpoint detail info
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        dsId path    string  true  "Data source id" gg
// @Param        id path    string  true  "Endpoint id" gg
// @Security bearerAuth
// @Success      200  {object}   DataSource
//
// @Router       /admin/ds/dse/{dsId}/{id} [get]
func dseItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(DataSourceEndpoint{}, w, r)
}

// deleteDs godoc
// @Summary      Delete data source
// @Description  Delete data source
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Ds id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/ds/{id} [delete]
func deleteDs(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(DataSource{}, w, r)
}

// deleteDse godoc
// @Summary      Delete data source endpoint
// @Description  Delete data source endpoint
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Data source endpoint id" id
// @Param        dsId    path     string  true  "Data source id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /ds/dse/{dsId}/{id} [delete]
func deleteDse(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(DataSource{}, w, r)
}

// updateDs
// @Summary      Update date source
// @Description  Update date source
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        source    body     DataSource  true  "Source info" true
// @Param        id    path     string  true  "Source info" id
// @Success      200 {object} DataSource
// @Security bearerAuth
//
// @Router       /admin/ds/{id} [put]
func updateDs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	exist, err := DataSource{}.GetById(vars["id"])

	if err != nil {
		w.WriteHeader(404)
		return
	}

	newm := DataSource{}

	err = json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.(DataSource).CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// updateDse
// @Summary      Update date source endpoint
// @Description  Update date source endpoint
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        dse    body     DataSourceEndpoint  true  "Endpoint info" true
// @Param        id    path     string  true  "Endpoint id" id
// @Param        dsId    path     string  true  "Data source id" id
// @Success      200 {object} DataSourceEndpoint
// @Security bearerAuth
//
// @Router       /admin/ds/dse/{dsId}/{id} [put]
func updateDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	exist, err := DataSourceEndpoint{}.GetById(vars["id"])

	if err != nil {
		w.WriteHeader(404)
		return
	}

	newm := DataSourceEndpoint{}

	err = json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.(DataSourceEndpoint).CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// publicDesItem godoc
// @Summary      Get item
// @Description  Get data source by id
// @Tags         Data source
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        db-key    header     string  true  "Auth key" gg
// @Param        id    path     string  true  "Source id"
// @Success      200  {object}   project.Project
//
// @Router       /dse/{id} [get]
func publicDesItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	m, err := DataSourceEndpoint{}.GetById(vars["id"])

	if err != nil {
		w.WriteHeader(404)
		return
	}

	model := m.(DataSourceEndpoint)

	arr, _ := model.List(10, 0, "id", "ASC", make(map[string]string))
	total := model.Total()
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, _ := json.Marshal(arr)

	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
