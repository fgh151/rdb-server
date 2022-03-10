package web

import (
	"db-server/meta"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ListTopics(w http.ResponseWriter, r *http.Request) {

	dpr := meta.MetaDb.List()

	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("Content-Type", "application/json")

	w.Header().Add("X-Total-Count", strconv.Itoa(len(dpr)))

	resp, _ := json.Marshal(dpr)
	w.WriteHeader(200)
	w.Write(resp)

}

func TopicItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resp, _ := json.Marshal(meta.MetaDb.GetById(vars["id"]))
	w.WriteHeader(200)
	w.Write(resp)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	meta.MetaDb.DeleteById(vars["id"])

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
