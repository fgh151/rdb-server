package utils

import (
	err2 "db-server/err"
	"db-server/modules"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func FormatQuery(r *http.Request, params []string) map[string]interface{} {
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

func ListItems(model modules.Model, filter []string, r *http.Request, w http.ResponseWriter) {
	log.Debug(r.Method, r.RequestURI)

	l, o, or, so := GetPagination(r)
	f := FormatQuery(r, filter)

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

func GetItem(m modules.Model, w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	resp, _ := json.Marshal(m.GetById(vars["id"]))
	w.WriteHeader(200)
	_, err := w.Write(resp)
	err2.DebugErr(err)
}

func DeleteItem(m modules.Model, w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	m.Delete(vars["id"])

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNoContent)
}
