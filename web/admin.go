package web

import (
	err2 "db-server/err"
	"db-server/meta"
	"db-server/models"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func getPagination(r *http.Request) (int, int, string, string) {

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

func listItems(model models.Model, r *http.Request, w http.ResponseWriter) {
	log.Debug(r.Method, r.RequestURI)
	arr := model.List(getPagination(r))
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

func CfLog(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	f := models.CloudFunction{}.GetById(vars["id"]).(models.CloudFunction)

	l, o, s, or := getPagination(r)
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

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = models.Project{}.GetById(vars["id"]).(models.Project)

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	meta.MetaDb.GetConnection().Save(&t)

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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	updateItem(models.User{}, w, r)
}

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	updateItem(models.Config{}, w, r)
}

func UpdateDs(w http.ResponseWriter, r *http.Request) {
	updateItem(models.DataSource{}, w, r)
}

func UpdateCf(w http.ResponseWriter, r *http.Request) {
	updateItem(models.CloudFunction{}, w, r)
}

func updateItem(m models.Model, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var t = m.GetById(vars["id"]).(models.Config)

	log.Debug(r.Method, r.RequestURI)
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	meta.MetaDb.GetConnection().Save(&t)

	resp, _ := json.Marshal(t)
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

	meta.MetaDb.GetConnection().Create(&t)

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
	meta.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

}

func CreateDs(w http.ResponseWriter, r *http.Request) {
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
	meta.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

func CreateCf(w http.ResponseWriter, r *http.Request) {
	model := models.CloudFunction{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	model.Id, err = uuid.NewUUID()
	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	meta.MetaDb.GetConnection().Create(&model)

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
