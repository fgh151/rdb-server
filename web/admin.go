package web

import (
	"db-server/meta"
	"db-server/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ListTopics(w http.ResponseWriter, r *http.Request) {
	listItems(models.Project{}.List(), w)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	listItems(models.User{}.List(), w)
}

func listItems(arr []interface{}, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")

	w.Header().Add("X-Total-Count", strconv.Itoa(len(arr)))

	resp, _ := json.Marshal(arr)
	w.WriteHeader(200)
	w.Write(resp)
}

func UserItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.User{}, w, r)
}

func TopicItem(w http.ResponseWriter, r *http.Request) {
	getItem(models.Project{}, w, r)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {

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
	w.Write(resp)
}

func getItem(m models.Model, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp, _ := json.Marshal(m.GetById(vars["id"]))
	w.WriteHeader(200)
	w.Write(resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.User{}, w, r)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	w.Write(resp)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	deleteItem(models.Project{}, w, r)
}

func deleteItem(m models.Model, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	m.Delete(vars["id"])

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {

	var t models.Project

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	meta.MetaDb.GetConnection().Create(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	w.Write(resp)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var t models.CreateUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user := t.Save()
	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	w.Write(resp)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var l models.LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := l.Login()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	w.Write(resp)

	return
}
