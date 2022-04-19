package web

import (
	"db-server/auth"
	"db-server/models"
	"encoding/json"
	"net/http"
)

func ApiAuth(w http.ResponseWriter, r *http.Request) {
	var l models.LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := l.ApiLogin()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	w.Write(resp)

	return
}

func ApiRegister(w http.ResponseWriter, r *http.Request) {

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

func ApiMe(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromRequest(r)
	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	w.Write(resp)
}
