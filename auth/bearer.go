package auth

import (
	"db-server/meta"
	"db-server/models"
	"encoding/json"
	"net/http"
	"strings"
)

func BearerVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := GetUserFromRequest(r)

		if user == nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Wrong auth token")
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
	meta.MetaDb.GetConnection().Find(&user, "token = ? ", reqToken)

	return user
}
