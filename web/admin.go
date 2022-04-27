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
	log.Debug(r.Method, r.RequestURI)
	listItems(models.Project{}, r, w)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	listItems(models.User{}, r, w)
}

func ListConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	listItems(models.Config{}, r, w)
}

func listItems(model models.Model, r *http.Request, w http.ResponseWriter) {
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
	log.Debug(r.Method, r.RequestURI)
	getItem(models.User{}, w, r)
}

func ConfigItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	getItem(models.Config{}, w, r)
}

func TopicItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
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
	log.Debug(r.Method, r.RequestURI)
	deleteItem(models.User{}, w, r)
}

func DeleteConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	deleteItem(models.Config{}, w, r)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = models.User{}.GetById(vars["id"]).(models.User)

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

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = models.Config{}.GetById(vars["id"]).(models.Config)

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
	log.Debug(r.Method, r.RequestURI)
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
	log.Debug(r.Method, r.RequestURI)

	var t models.Config
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
