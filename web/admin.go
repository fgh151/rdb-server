package web

import (
	"db-server/meta"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ListTopics(w http.ResponseWriter, r *http.Request) {

	dpr := meta.MetaDb.ListProjects()

	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")

	w.Header().Add("X-Total-Count", strconv.Itoa(len(dpr)))

	resp, _ := json.Marshal(dpr)
	w.WriteHeader(200)
	w.Write(resp)

}

func ListUsers(w http.ResponseWriter, r *http.Request) {

	dpr := meta.MetaDb.ListUsers()

	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")

	w.Header().Add("X-Total-Count", strconv.Itoa(len(dpr)))

	resp, _ := json.Marshal(dpr)
	w.WriteHeader(200)
	w.Write(resp)

}

type createUserForm struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var t createUserForm

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var u = meta.User{
		Email:        t.Email,
		PasswordHash: meta.HashPassword(t.Password),
		Token:        meta.GenerateRandomString(15),
	}

	meta.MetaDb.GetConnection().Create(&u)

	resp, _ := json.Marshal(u)
	w.WriteHeader(200)
	w.Write(resp)
}

func UserItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp, _ := json.Marshal(meta.MetaDb.GetUserById(vars["id"]))
	w.WriteHeader(200)
	w.Write(resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	meta.MetaDb.DeleteUserById(vars["id"])

	w.WriteHeader(200)
}

func TopicItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp, _ := json.Marshal(meta.MetaDb.GetProjectById(vars["id"]))
	w.WriteHeader(200)
	w.Write(resp)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	meta.MetaDb.DeleteProjectById(vars["id"])

	w.WriteHeader(200)
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {

	var t meta.Project

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

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var l LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var user meta.User

	meta.MetaDb.GetConnection().First(&user, "email = ?", l.Email)

	if user.ValidatePassword(l.Password) {
		resp, _ := json.Marshal(user)
		w.WriteHeader(200)
		w.Write(resp)
	}

	return
}
