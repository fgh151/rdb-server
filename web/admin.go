package web

import (
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/models"
	"db-server/server"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetPagination(r *http.Request) (int, int, string, string) {

	v := r.URL.Query()

	limit, err := strconv.Atoi(v.Get("_end"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(v.Get("_start"))
	if err != nil {
		offset = 0
	}

	order := v.Get("_order")
	if order == "" {
		order = "ASC"
	}
	sort := v.Get("_sort")
	if sort == "" {
		sort = "id"
	}

	return limit - offset, offset, order, sort
}

func formatQuery(r *http.Request, params []string) map[string]interface{} {
	result := make(map[string]interface{})

	if len(params) < 1 {
		return result
	}

	v := r.URL.Query()
	for _, param := range params {
		if v.Has(param) {
			val := v.Get(param)
			if val != "" {
				result[param] = val
			}
		}
	}

	return result
}

// ListTopics godoc
// @Summary      List topics
// @Description  List topics
// @Tags         TopicOutput
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.Project
//
// @Router       /admin/topics [get]
func ListTopics(w http.ResponseWriter, r *http.Request) {
	listItems(models.Project{}, []string{}, r, w)
}

// ListUsers godoc
// @Summary      List users
// @Description  List users
// @Tags         User
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.User
//
// @Router       /admin/users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	listItems(models.User{}, []string{"id", "email", "admin", "active"}, r, w)
}

// ListConfig godoc
// @Summary      List configs
// @Description  List configs
// @Tags         Config manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.Config
//
// @Router       /admin/config [get]
func ListConfig(w http.ResponseWriter, r *http.Request) {
	listItems(models.Config{}, []string{}, r, w)
}

// ListDs godoc
// @Summary      List data sources
// @Description  List data sources
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.DataSource
//
// @Router       /admin/ds [get]
func ListDs(w http.ResponseWriter, r *http.Request) {
	listItems(models.DataSource{}, []string{}, r, w)
}

// ListDse godoc
// @Summary      List data source endpoints
// @Description  List data source endpoints
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Param        dsIid path    string  true  "Data source id" id
// @Success      200  {array}   models.DataSourceEndpoint
//
// @Router       /admin/ds/{dsIid}/dse [get]
func ListDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	l, o, or, so := GetPagination(r)
	f := formatQuery(r, []string{"data_source_id"})

	arr := models.DataSourceEndpoint{}.List(l, o, so, or, f)
	total := len(arr)
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", strconv.Itoa(total))

	resp, _ := json.Marshal(arr)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// ListCf godoc
// @Summary      List cloud functions
// @Description  List cloud functions
// @Tags         Cloud functions
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.CloudFunction
//
// @Router       /admin/cf [get]
func ListCf(w http.ResponseWriter, r *http.Request) {
	listItems(models.CloudFunction{}, []string{"id"}, r, w)
}

// ListPipeline godoc
// @Summary      List pipelines
// @Description  List pipelines
// @Tags         Pipeline
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.Pipeline
//
// @Router       /admin/pl [get]
func ListPipeline(w http.ResponseWriter, r *http.Request) {
	listItems(models.Pipeline{}, []string{"id"}, r, w)
}

// ListPush godoc
// @Summary      List push messages
// @Description  List push messages
// @Tags         Push messages
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.PushMessage
//
// @Router       /admin/push [get]
func ListPush(w http.ResponseWriter, r *http.Request) {
	listItems(models.PushMessage{}, []string{}, r, w)
}

// ListCron godoc
// @Summary      List cron jobs
// @Description  List cron jobs
// @Tags         Cron
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.CronJob
//
// @Router       /admin/cron [get]
func ListCron(w http.ResponseWriter, r *http.Request) {
	listItems(models.CronJob{}, []string{}, r, w)
}

func listItems(model models.Model, filter []string, r *http.Request, w http.ResponseWriter) {
	log.Debug(r.Method, r.RequestURI)

	l, o, or, so := GetPagination(r)
	f := formatQuery(r, filter)

	arr := model.List(l, o, so, or, f)
	total := model.Total()
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, err := json.Marshal(arr)
	log.Debug(string(resp))
	err2.DebugErr(err)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// SettingsOauth godoc
// @Summary      List OAuth settings
// @Description  List OAuth settings
// @Tags         OAuth
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.AppSettings
//
// @Router       /settings/{projectId}/oauth [get]
func SettingsOauth(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	pid := params["projectId"]

	var result []map[string]interface{}
	tx := server.MetaDb.GetConnection().Table("app_settings").Where("name like ? AND project_id = ?", "oauth_%", pid).Select("name", "value").Find(&result)

	log.Debug(tx.RowsAffected)

	resp, err := json.Marshal(result)
	err2.DebugErr(err)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UserItem godoc
// @Summary      User info
// @Description  User detail info
// @Tags         User
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "User id" gg
// @Security bearerAuth
// @Success      200  {object}   models.User
//
// @Router       /admin/users/{id} [get]
func UserItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.User{}, w, r)
}

// ConfigItem godoc
// @Summary      Config info
// @Description  Config detail info
// @Tags         Config manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Config id" gg
// @Security bearerAuth
// @Success      200  {object}   models.Config
//
// @Router       /admin/config/{id} [get]
func ConfigItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.Config{}, w, r)
}

// DsItem godoc
// @Summary      Data source info
// @Description  Data source detail info
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Source id" gg
// @Security bearerAuth
// @Success      200  {object}   models.DataSource
//
// @Router       /admin/ds/{id} [get]
func DsItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.DataSource{}, w, r)
}

// DseItem godoc
// @Summary      Data source endpoint info
// @Description  Data source endpoint detail info
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        dsId path    string  true  "Data source id" gg
// @Param        id path    string  true  "Endpoint id" gg
// @Security bearerAuth
// @Success      200  {object}   models.DataSource
//
// @Router       /admin/ds/dse/{dsId}/{id} [get]
func DseItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.DataSourceEndpoint{}, w, r)
}

// CfItem godoc
// @Summary      Function info
// @Description  Function detail info
// @Tags         Cloud functions
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "cf id" id
// @Security bearerAuth
// @Success      200  {object}   models.CloudFunction
//
// @Router       /admin/cf/{id} [get]
func CfItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.CloudFunction{}, w, r)
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
// @Success      200  {object}   models.Pipeline
//
// @Router       /admin/pl/{id} [get]
func PipelineItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.Pipeline{}, w, r)
}

// PushItem godoc
// @Summary      Push info
// @Description  Push detail info
// @Tags         Push messages
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Push id" id
// @Security bearerAuth
// @Success      200  {object}   models.PushMessage
//
// @Router       /admin/push/{id} [get]
func PushItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.PushMessage{}, w, r)
}

// CronItem godoc
// @Summary      Cron job info
// @Description  Cron job detail info
// @Tags         Cron
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Push id" id
// @Security bearerAuth
// @Success      200  {object}   models.CronJob
//
// @Router       /admin/cron/{id} [get]
func CronItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.CronJob{}, w, r)
}

// CfLog godoc
// @Summary      Logs
// @Description  Cloud function logs
// @Tags         Cloud functions
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Fuc id" id
// @Security bearerAuth
// @Success      200  {object}   models.CloudFunctionLog
//
// @Router       /admin/cf/{id}/log [get]
func CfLog(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	f := models.CloudFunction{}.GetById(vars["id"]).(models.CloudFunction)

	l, o, s, or := GetPagination(r)
	arr := models.ListCfLog(f.Id, l, o, s, or)
	total := models.LogsTotal(f.Id)
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, _ := json.Marshal(arr)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// TopicItem godoc
// @Summary      TopicOutput
// @Description  topic detail info
// @Tags         Entity manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "TopicOutput id" id
// @Security bearerAuth
// @Success      200  {object}   models.Project
//
// @Router       /admin/topics/{id} [get]
func TopicItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.Project{}, w, r)
}

// TopicData godoc
// @Summary      TopicOutput data
// @Description  topic data
// @Tags         Entity manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        topic path    string  true  "TopicOutput name"
// @Security bearerAuth
// @Success      200  {array} object
//
// @Router       /admin/topics/{topic}/data [get]
func TopicData(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	limit, offset, rorder, sort := GetPagination(r)

	order, sort := drivers.GetMongoSort(sort, rorder)

	log.Debug("Mongo limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " order " + rorder + " sort " + sort)

	res, count, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic, int64(limit), int64(offset), order, sort, bson.D{})

	var result []map[string]string

	for _, resArray := range res {
		record := make(map[string]string)
		for key, obj := range resArray.Map() {

			if key == "_id" {
				key = "id"
				obj = strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", obj), "bjectID(\"", ""), "\")", "")
			}
			record[key] = fmt.Sprintf("%v", obj)
		}
		result = append(result, record)
	}

	w.Header().Add("X-Total-Count", strconv.FormatInt(count, 10))

	sendResponse(w, 200, result, err)
}

// UpdateTopic
// @Summary      Update topic
// @Description  Update topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.Project  true  "Project info" true
// @Param        id    path     string  true  "Project id" true
// @Success      200 {object} models.Project
// @Security bearerAuth
//
// @Router       /admin/topics/{id} [put]
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = models.Project{}.GetById(vars["id"]).(models.Project)

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	server.MetaDb.GetConnection().Save(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

func getItem(m models.Model, w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	resp, _ := json.Marshal(m.GetById(vars["id"]))
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

// DeleteUser godoc
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
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.User{}, w, r)
}

// DeleteConfig godoc
// @Summary      Delete config
// @Description  Delete config
// @Tags         Config manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Config id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/config/{id} [delete]
func DeleteConfig(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.Config{}, w, r)
}

// DeleteDs godoc
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
func DeleteDs(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.DataSource{}, w, r)
}

// DeleteDse godoc
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
func DeleteDse(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.DataSource{}, w, r)
}

// DeleteCf godoc
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
func DeleteCf(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.CloudFunction{}, w, r)
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
	deleteItem(models.Pipeline{}, w, r)
}

// DeletePush godoc
// @Summary      Delete push
// @Description  Delete push
// @Tags         Push messages
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Push id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/push/{id} [delete]
func DeletePush(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.PushMessage{}, w, r)
}

// DeleteCron godoc
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
func DeleteCron(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.CronJob{}, w, r)
}

// UpdateUser
// @Summary      Update user
// @Description  Update user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.User  true  "User info" true
// @Param        id    path     string  true  "User info" id
// @Success      200 {object} models.User
// @Security bearerAuth
//
// @Router       /admin/users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.User{}.GetById(vars["id"]).(models.User)
	newm := models.User{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt
	newm.LastLogin = exist.LastLogin
	newm.PasswordHash = exist.PasswordHash

	server.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UpdateConfig
// @Summary      Update config
// @Description  Update config
// @Tags         Config manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.Config  true  "Config info" true
// @Param        id    path     string  true  "Config id" id
// @Success      200 {object} models.Config
// @Security bearerAuth
//
// @Router       /admin/config/{id} [put]
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	newm := models.Config{}

	err := json.NewDecoder(r.Body).Decode(&newm)
	server.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UpdateDs
// @Summary      Update date source
// @Description  Update date source
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        source    body     models.DataSource  true  "Source info" true
// @Param        id    path     string  true  "Source info" id
// @Success      200 {object} models.DataSource
// @Security bearerAuth
//
// @Router       /admin/ds/{id} [put]
func UpdateDs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.DataSource{}.GetById(vars["id"]).(models.DataSource)
	newm := models.DataSource{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	server.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UpdateDse
// @Summary      Update date source endpoint
// @Description  Update date source endpoint
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        dse    body     models.DataSourceEndpoint  true  "Endpoint info" true
// @Param        id    path     string  true  "Endpoint id" id
// @Param        dsId    path     string  true  "Data source id" id
// @Success      200 {object} models.DataSourceEndpoint
// @Security bearerAuth
//
// @Router       /admin/ds/dse/{dsId}/{id} [put]
func UpdateDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.DataSourceEndpoint{}.GetById(vars["id"]).(models.DataSourceEndpoint)
	newm := models.DataSourceEndpoint{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	server.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UpdateCf
// @Summary      Update cloud function
// @Description  Update cloud function
// @Tags         Cloud functions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.DataSource  true  "Source info" true
// @Param        id    path     string  true  "Function id" id
// @Success      200 {object} models.DataSource
// @Security bearerAuth
//
// @Router       /admin/cf/{id} [put]
func UpdateCf(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)

	var projectId, _ = uuid.Parse(r.FormValue("project_id"))

	uri, err := models.GetContainerUri(r.FormValue("container"))
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

	server.MetaDb.GetConnection().Table("cloud_functions").Where("id = ?", vars["id"]).Updates(
		map[string]interface{}{
			"title":      r.FormValue("title"),
			"project_id": projectId,
			"container":  r.FormValue("container"),
			"params":     r.FormValue("params"),
			"env":        r.FormValue("env"),
		},
	)

	resp, _ := json.Marshal(models.CloudFunction{}.GetById(vars["id"]))
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

}

// UpdatePipeline
// @Summary      Update pipeline
// @Description  Update pipeline
// @Tags         Pipeline
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.Pipeline  true  "Pipeline info" true
// @Param        id    path     string  true  "Pipeline id" id
// @Success      200 {object} models.Pipeline
// @Security bearerAuth
//
// @Router       /admin/pl/{id} [put]
func UpdatePipeline(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.Pipeline{}.GetById(vars["id"]).(models.Pipeline)
	newm := models.Pipeline{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	server.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UpdatePush
// @Summary      Update push
// @Description  Update push
// @Tags         Push messages
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.PushMessage  true  "push info" true
// @Param        id    path     string  true  "push id" id
// @Success      200 {object} models.PushMessage
// @Security bearerAuth
//
// @Router       /admin/push/{id} [put]
func UpdatePush(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.PushMessage{}.GetById(vars["id"]).(models.PushMessage)
	newm := models.PushMessage{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	server.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// UpdateCron
// @Summary      Update cron job
// @Description  Update cron job
// @Tags         Cron
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.CronJob  true  "Cron job info" true
// @Param        id    path     string  true  "Cron id"
// @Success      200 {object} models.CronJob
// @Security bearerAuth
//
// @Router       /admin/cron/{id} [put]
func UpdateCron(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.CronJob{}.GetById(vars["id"]).(models.CronJob)
	newm := models.CronJob{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt
	server.MetaDb.GetConnection().Save(&newm)

	c := server.Cron
	c.GetScheduler().Remove(exist.CronId)
	newm.Schedule(c.GetScheduler())

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// DeleteTopic godoc
// @Summary      Delete topic
// @Description  Delete topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "TopicOutput id" string
// @Success      204
//
// @Router       /admin/topics/{id} [delete]
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.Project{}, w, r)
}

func deleteItem(m models.Model, w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	m.Delete(vars["id"])

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
}

// CreateTopic
// @Summary      Create topic
// @Description  Create topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        topic    body     models.Project  true  "topic info" true
// @Success      200 {object} models.Project
// @Security bearerAuth
//
// @Router       /admin/topics [post]
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t models.Project
	t.Id, _ = uuid.NewUUID()

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	server.MetaDb.GetConnection().Create(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreateUser
// @Summary      Create user
// @Description  Create user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        user    body     models.CreateUserForm  true  "User info" true
// @Success      200 {object} models.User
// @Security bearerAuth
//
// @Router       /admin/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t models.CreateUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user := t.Save()
	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreateConfig
// @Summary      Create config
// @Description  Create config
// @Tags         Config manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        config    body     models.CreateUserForm  true  "Config info" true
// @Success      200 {object} models.User
// @Security bearerAuth
//
// @Router       /admin/config [post]
func CreateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.Config{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	model.Id, err = uuid.NewUUID()
	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

}

// CreateDs
// @Summary      Create data source
// @Description  Create data source
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        ds    body     models.DataSource  true  "Data source info" true
// @Success      200 {object} models.DataSource
// @Security bearerAuth
//
// @Router       /admin/ds [post]
func CreateDs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.DataSource{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	model.Id, err = uuid.NewUUID()
	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreateDse
// @Summary      Create data source endpoint
// @Description  Create data source endpoint
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        dse    body     models.DataSourceEndpoint  true  "Data source info" true
// @Param        dsId    path     string  true  "Data source id" id
// @Success      200 {object} models.DataSourceEndpoint
// @Security bearerAuth
//
// @Router       /admin/ds/dse/{dsId} [post]
func CreateDse(w http.ResponseWriter, r *http.Request) {
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

	model := models.DataSourceEndpoint{
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
	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreateCf
// @Summary      Create cloud function
// @Description  Create cloud function
// @Tags         Cloud functions
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        cf    body     models.CloudFunction  true  "Function info" true
// @Success      200 {object} models.CloudFunction
// @Security bearerAuth
//
// @Router       /admin/cf [post]
func CreateCf(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.CloudFunction{}

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
		uri, err := models.GetContainerUri(model.Container)
		err2.DebugErr(err)

		go func() {
			err := server.BuildDockerImage(file, []string{uri.Vendor + "/" + uri.Image})
			err2.DebugErr(err)
		}()
	}

	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreatePipeline
// @Summary      Create pipeline
// @Description  Create pipeline
// @Tags         Pipeline
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        pl    body     models.Pipeline  true  "Pipeline info" true
// @Success      200 {object} models.Pipeline
// @Security bearerAuth
//
// @Router       /admin/pl [post]
func CreatePipeline(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.Pipeline{}

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

	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreatePush
// @Summary      Create push message
// @Description  Create push message
// @Tags         Push messages
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        push    body     models.PushMessage  true  "Push info" true
// @Success      200 {object} models.PushMessage
// @Security bearerAuth
//
// @Router       /admin/push [post]
func CreatePush(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.PushMessage{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	id, err := uuid.NewUUID()
	model.Id = id
	model.Sent = false

	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreateCron
// @Summary      Create cron job
// @Description  Create cron job
// @Tags         Cron
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        cron    body     models.CronJob  true  "Push info" true
// @Success      200 {object} models.CronJob
// @Security bearerAuth
//
// @Router       /admin/cron [post]
func CreateCron(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.CronJob{}

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
	server.MetaDb.GetConnection().Create(&model)

	model.Schedule(server.Cron.GetScheduler())

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// Auth godoc
// @Summary      Login to admin
// @Description  Authenticate in admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   models.User
//
// @Router       /admin/auth [post]
func Auth(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var l models.LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := l.AdminLogin()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

	return
}
