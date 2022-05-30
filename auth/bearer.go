package auth

import (
	err2 "db-server/err"
	"db-server/models"
	"db-server/server"
	"encoding/json"
	"net/http"
	"strings"
)

func BearerVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := GetUserFromRequest(r)

		if user == nil {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode("Wrong auth token")
			err2.DebugErr(err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserFromRequest(r *http.Request) *models.User {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		return nil
	}

	reqToken = splitToken[1]

	var user *models.User
	server.MetaDb.GetConnection().Find(&user, "token = ? ", reqToken)

	return user
}
