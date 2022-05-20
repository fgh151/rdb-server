package web

import (
	"db-server/auth"
	err2 "db-server/err"
	"db-server/models"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func ApiAuth(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var l models.LoginForm
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := l.ApiLogin()

	if err != nil {
		err2.DebugErr(err)
		http.Error(w, err.Error(), 400)
		return
	}

	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)

	return
}

func ApiRegister(w http.ResponseWriter, r *http.Request) {
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

func ApiMe(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	user := auth.GetUserFromRequest(r)
	resp, _ := json.Marshal(user)
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

func ApiConfigItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	model := models.Config{}.GetById(vars["id"]).(models.Config)

	rKey := r.Header.Get("db-key")

	if !validateKey(model.Project.Key, rKey) {
		send403Error(w, "db-key not Valid")
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, err := w.Write([]byte(model.Body))
		err2.DebugErr(err)
	}
}

func DSEItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	model := models.DataSourceEndpoint{}.GetById(vars["id"]).(models.DataSourceEndpoint)

	arr := model.List(10, 0, "id", "ASC")
	total := model.Total()
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("X-Total-Count", strconv.FormatInt(*total, 10))

	resp, _ := json.Marshal(arr)

	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

func CfRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cf := models.CloudFunction{}.GetById(vars["id"]).(models.CloudFunction)
	cf.Run()
	w.WriteHeader(200)
}
