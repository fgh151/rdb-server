package utils

import (
	err2 "db-server/err"
	"db-server/server/db"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
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

	order := v.Get("_order")
	if order == "" {
		order = "ASC"
	}
	sort := v.Get("_sort")
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

// GetUserFromRequest Fetch user model from request
func GetUserFromRequest(r *http.Request) (interface{}, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		return nil, nil
	}

	reqToken = splitToken[1]

	var usr *interface{}
	tx := db.MetaDb.GetConnection().Find(&usr, "token = ? ", reqToken)

	if tx.RowsAffected < 1 {
		return usr, errors.New("invalid credentials")
	}

	return usr, nil
}

func GetPayload(r *http.Request) map[string]interface{} {
	var requestPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPayload)
	err2.PanicErr(err)

	return requestPayload
}
