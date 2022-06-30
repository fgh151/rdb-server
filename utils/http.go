package utils

import (
	err2 "db-server/err"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

	order := CleanInputString(v.Get("_order"))
	if order == "" {
		order = "ASC"
	}
	sort := CleanInputString(v.Get("_sort"))
	if sort == "" {
		sort = "id"
	}

	return limit - offset, offset, order, sort
}

func Send403Error(w http.ResponseWriter, message string) {
	logrus.Debug("403 error")
	payload := map[string]string{"code": "not acceptable", "message": message}
	SendResponse(w, 403, payload, nil)
}

func SendResponse(w http.ResponseWriter, statusCode int, payload interface{}, err error) {
	if err == nil {
		w.WriteHeader(statusCode)
		if payload != nil {
			resp, _ := json.Marshal(payload)
			_, err = w.Write(resp)
			err2.DebugErr(err)
		}
	} else {
		w.WriteHeader(500)
		_, err = w.Write([]byte(err.Error()))
		err2.DebugErr(err)
	}
}

func GetPayload(r *http.Request) map[string]interface{} {
	var requestPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPayload)
	err2.PanicErr(err)

	return requestPayload
}
