package web

import (
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/modules/cf"
	"db-server/modules/config"
	"db-server/modules/cron"
	"db-server/modules/ds"
	"db-server/modules/pipeline"
	"db-server/modules/project"
	"db-server/modules/user"
	"db-server/server"
	"db-server/server/db"
	"db-server/utils"
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

// ListTopics godoc
// @Summary      List topics
// @Description  List topics
// @Tags         TopicOutput
// @Accept       json
// @Produce      json
// @Security bearerAuth
/* @Success      200  {array}   project.Project */
//
// @Router       /admin/topics [get]
func ListTopics(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(project.Project{}, []string{}, r, w)
}

// ListUsers godoc
// @Summary      List users
// @Description  List users
// @Tags         User
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   user.User
//
// @Router       /admin/users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(user.User{}, []string{"id", "email", "admin", "active"}, r, w)
}

// ListConfig godoc
// @Summary      List configs
// @Description  List configs
// @Tags         Config manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   config.Config
//
// @Router       /admin/config [get]
func ListConfig(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(config.Config{}, []string{}, r, w)
}

// ListDs godoc
// @Summary      List data sources
// @Description  List data sources
// @Tags         Data source
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   ds.DataSource
//
// @Router       /admin/ds [get]
func ListDs(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(ds.DataSource{}, []string{}, r, w)
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
// @Success      200  {array}   ds.DataSourceEndpoint
//
// @Router       /admin/ds/{dsIid}/dse [get]
func ListDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	l, o, or, so := utils.GetPagination(r)
	f := utils.FormatQuery(r, []string{"data_source_id"})

	arr := ds.DataSourceEndpoint{}.List(l, o, so, or, f)
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
// @Success      200  {array}   cf.CloudFunction
//
// @Router       /admin/cf [get]
func ListCf(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(cf.CloudFunction{}, []string{"id"}, r, w)
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
	utils.ListItems(pipeline.Pipeline{}, []string{"id"}, r, w)
}

// ListCron godoc
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
func ListCron(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(cron.CronJob{}, []string{}, r, w)
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
// @Success      200  {object}   user.User
//
// @Router       /admin/users/{id} [get]
func UserItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(user.User{}, w, r)
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
// @Success      200  {object}   config.Config
//
// @Router       /admin/config/{id} [get]
func ConfigItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(config.Config{}, w, r)
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
// @Success      200  {object}   ds.DataSource
//
// @Router       /admin/ds/{id} [get]
func DsItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(ds.DataSource{}, w, r)
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
// @Success      200  {object}   ds.DataSource
//
// @Router       /admin/ds/dse/{dsId}/{id} [get]
func DseItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(ds.DataSourceEndpoint{}, w, r)
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
// @Success      200  {object}   cf.CloudFunction
//
// @Router       /admin/cf/{id} [get]
func CfItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(cf.CloudFunction{}, w, r)
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
	utils.GetItem(pipeline.Pipeline{}, w, r)
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
// @Success      200  {object}   cron.CronJob
//
// @Router       /admin/cron/{id} [get]
func CronItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(cron.CronJob{}, w, r)
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
// @Success      200  {object}   cf.CloudFunctionLog
//
// @Router       /admin/cf/{id}/log [get]
func CfLog(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	f := cf.CloudFunction{}.GetById(vars["id"]).(cf.CloudFunction)

	l, o, s, or := utils.GetPagination(r)
	arr := cf.ListCfLog(f.Id, l, o, s, or)
	total := cf.LogsTotal(f.Id)
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
/* @Success      200  {object}   project.Project */
//
// @Router       /admin/topics/{id} [get]
func TopicItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(project.Project{}, w, r)
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

	limit, offset, rorder, sort := utils.GetPagination(r)

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
/* / @Param        device    body     project.Project  true  "Project info" true
// @Param        id    path     string  true  "Project id" true
// @Success      200 {object} project.Project */
// @Security bearerAuth
//
// @Router       /admin/topics/{id} [put]
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = project.Project{}.GetById(vars["id"]).(project.Project)

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
	utils.DeleteItem(user.User{}, w, r)
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
	utils.DeleteItem(config.Config{}, w, r)
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
	utils.DeleteItem(ds.DataSource{}, w, r)
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
	utils.DeleteItem(ds.DataSource{}, w, r)
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
	utils.DeleteItem(cf.CloudFunction{}, w, r)
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
	utils.DeleteItem(pipeline.Pipeline{}, w, r)
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
	utils.DeleteItem(cron.CronJob{}, w, r)
}

// UpdateUser
// @Summary      Update user
// @Description  Update user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     user.User  true  "User info" true
// @Param        id    path     string  true  "User info" id
// @Success      200 {object} user.User
// @Security bearerAuth
//
// @Router       /admin/users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = user.User{}.GetById(vars["id"]).(user.User)
	newm := user.User{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt
	newm.LastLogin = exist.LastLogin
	newm.PasswordHash = exist.PasswordHash

	db.MetaDb.GetConnection().Save(&newm)

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
// @Param        device    body     config.Config  true  "Config info" true
// @Param        id    path     string  true  "Config id" id
// @Success      200 {object} config.Config
// @Security bearerAuth
//
// @Router       /admin/config/{id} [put]
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	newm := config.Config{}

	err := json.NewDecoder(r.Body).Decode(&newm)
	db.MetaDb.GetConnection().Save(&newm)

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
// @Param        source    body     ds.DataSource  true  "Source info" true
// @Param        id    path     string  true  "Source info" id
// @Success      200 {object} ds.DataSource
// @Security bearerAuth
//
// @Router       /admin/ds/{id} [put]
func UpdateDs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = ds.DataSource{}.GetById(vars["id"]).(ds.DataSource)
	newm := ds.DataSource{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

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
// @Param        dse    body     ds.DataSourceEndpoint  true  "Endpoint info" true
// @Param        id    path     string  true  "Endpoint id" id
// @Param        dsId    path     string  true  "Data source id" id
// @Success      200 {object} ds.DataSourceEndpoint
// @Security bearerAuth
//
// @Router       /admin/ds/dse/{dsId}/{id} [put]
func UpdateDse(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = ds.DataSourceEndpoint{}.GetById(vars["id"]).(ds.DataSourceEndpoint)
	newm := ds.DataSourceEndpoint{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

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
// @Param        device    body     ds.DataSource  true  "Source info" true
// @Param        id    path     string  true  "Function id" id
// @Success      200 {object} ds.DataSource
// @Security bearerAuth
//
// @Router       /admin/cf/{id} [put]
func UpdateCf(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)

	var projectId, _ = uuid.Parse(r.FormValue("project_id"))

	uri, err := cf.GetContainerUri(r.FormValue("container"))
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

	resp, _ := json.Marshal(cf.CloudFunction{}.GetById(vars["id"]))
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
// @Param        device    body     pipeline.Pipeline  true  "Pipeline info" true
// @Param        id    path     string  true  "Pipeline id" id
// @Success      200 {object} pipeline.Pipeline
// @Security bearerAuth
//
// @Router       /admin/pl/{id} [put]
func UpdatePipeline(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = pipeline.Pipeline{}.GetById(vars["id"]).(pipeline.Pipeline)
	newm := pipeline.Pipeline{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

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
// @Param        device    body     cron.CronJob  true  "Cron job info" true
// @Param        id    path     string  true  "Cron id"
// @Success      200 {object} cron.CronJob
// @Security bearerAuth
//
// @Router       /admin/cron/{id} [put]
func UpdateCron(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = cron.CronJob{}.GetById(vars["id"]).(cron.CronJob)
	newm := cron.CronJob{}

	err := json.NewDecoder(r.Body).Decode(&newm)

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
	utils.DeleteItem(project.Project{}, w, r)
}

// CreateTopic
// @Summary      Create topic
// @Description  Create topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
/* @Param        topic    body     project.Project  true  "topic info" true
// @Success      200 {object} project.Project */
// @Security bearerAuth
//
// @Router       /admin/topics [post]
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t project.Project
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

// CreateUser
// @Summary      Create user
// @Description  Create user
// @Tags         User
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        user    body     user.CreateUserForm  true  "User info" true
// @Success      200 {object} user.User
// @Security bearerAuth
//
// @Router       /admin/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t user.CreateUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	usr := t.Save()
	resp, _ := json.Marshal(usr)
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
// @Param        config    body     user.CreateUserForm  true  "Config info" true
// @Success      200 {object} user.User
// @Security bearerAuth
//
// @Router       /admin/config [post]
func CreateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := config.Config{}

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

// CreateDs
// @Summary      Create data source
// @Description  Create data source
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        ds    body     ds.DataSource  true  "Data source info" true
// @Success      200 {object} ds.DataSource
// @Security bearerAuth
//
// @Router       /admin/ds [post]
func CreateDs(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := ds.DataSource{}

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

// CreateDse
// @Summary      Create data source endpoint
// @Description  Create data source endpoint
// @Tags         Data source
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        dse    body     ds.DataSourceEndpoint  true  "Data source info" true
// @Param        dsId    path     string  true  "Data source id" id
// @Success      200 {object} ds.DataSourceEndpoint
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

	model := ds.DataSourceEndpoint{
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

// CreateCf
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
func CreateCf(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := cf.CloudFunction{}

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
		uri, err := cf.GetContainerUri(model.Container)
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
	model := pipeline.Pipeline{}

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

// CreateCron
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
func CreateCron(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := cron.CronJob{}

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

// Auth godoc
// @Summary      Login to admin
// @Description  Authenticate in admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        email    query     string  false  "Email for login" gg
// @Param        password    query     string  false  "Password for login" gg
// @Success      200  {object}   user.User
//
// @Router       /admin/auth [post]
func Auth(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var l user.LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	usr, err := l.AdminLogin()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(usr)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

	return
}
