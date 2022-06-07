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
		order = "id"
	}
	sort := v.Get("_sort")
	if sort == "" {
		sort = "ASC"
	}

	return limit - offset, offset, order, sort
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	listItems(models.Project{}, r, w)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	listItems(models.User{}, r, w)
}

func ListConfig(w http.ResponseWriter, r *http.Request) {
	listItems(models.Config{}, r, w)
}

func ListDs(w http.ResponseWriter, r *http.Request) {
	listItems(models.DataSource{}, r, w)
}

func ListCf(w http.ResponseWriter, r *http.Request) {
	listItems(models.CloudFunction{}, r, w)
}

func ListPush(w http.ResponseWriter, r *http.Request) {
	listItems(models.PushMessage{}, r, w)
}

func ListCron(w http.ResponseWriter, r *http.Request) {
	listItems(models.CronJob{}, r, w)
}

func listItems(model models.Model, r *http.Request, w http.ResponseWriter) {
	log.Debug(r.Method, r.RequestURI)
	arr := model.List(GetPagination(r))
	total := model.Total()
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, _ := json.Marshal(arr)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

func UserItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.User{}, w, r)
}

func ConfigItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.Config{}, w, r)
}

func DsItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.DataSource{}, w, r)
}

func CfItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.CloudFunction{}, w, r)
}

func PushItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.PushMessage{}, w, r)
}

func CronItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.CronJob{}, w, r)
}

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

func TopicItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.Project{}, w, r)
}

func TopicData(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	limit, offset, rorder, sort := GetPagination(r)

	order, sort := drivers.GetMongoSort(sort, rorder)

	log.Debug("Mongo limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " order " + rorder + " sort " + sort)

	res, count, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic, int64(limit), int64(offset), order, sort)

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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.User{}, w, r)
}

func DeleteConfig(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.Config{}, w, r)
}

func DeleteDs(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.DataSource{}, w, r)
}

func DeleteCf(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.CloudFunction{}, w, r)
}

func DeletePush(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.PushMessage{}, w, r)
}

func DeleteCron(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.CronJob{}, w, r)
}

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

func UpdateCf(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)

	var projectId, _ = uuid.Parse(r.FormValue("project_id"))

	file, _, err := r.FormFile("dockerarc")
	if err == nil {
		uri, err := models.GetContainerUri(r.FormValue("container"))
		err2.DebugErr(err)

		go func() {
			err := models.BuildImage(file, uri)
			err2.DebugErr(err)
		}()
	} else {
		log.Debug(err)
	}

	server.MetaDb.GetConnection().Table("cloud_functions").Where("id = ?", vars["id"]).Updates(
		map[string]interface{}{"title": r.FormValue("title"), "project_id": projectId, "container": r.FormValue("container"), "params": r.FormValue("params")},
	)

	w.WriteHeader(200)
	err2.DebugErr(err)
}

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
			err := models.BuildImage(file, uri)
			err2.DebugErr(err)
		}()
	}

	server.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

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
