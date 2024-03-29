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

func FormatQuery(r *http.Request, params []string) map[string]string {
	result := make(map[string]string)

	if len(params) < 1 {
		return result
	}

	v := r.URL.Query()
	for _, param := range params {
		if v.Has(param) {
			val := v.Get(param)
			if val != "" {
				result[CleanInputString(param)] = CleanInputString(val)
			}
		}
	}

	return result
}

func ListItems(model modules.Model, filter []string, r *http.Request, w http.ResponseWriter) {
	log.Debug(r.Method, r.RequestURI)

	l, o, or, so := GetPagination(r)
	f := FormatQuery(r, filter)

	arr, _ := model.List(l, o, so, or, f)
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

	model, err := m.GetById(vars["id"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
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
