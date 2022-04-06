package auth

import (
	"db-server/db"
	"db-server/models"
	"encoding/json"
	"net/http"
	"strings"
)

func BearerVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) < 2 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Missing auth token")
			return
		}

		reqToken = splitToken[1]

		var user *models.User
		db.DB.GetConnection().Find(&user, "token = ? ", reqToken)

		if user == nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Wrong auth token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
