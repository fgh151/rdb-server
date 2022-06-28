package auth

import (
	err2 "db-server/err"
	"db-server/modules/user"
	"db-server/server/db"
	"encoding/json"
	"net/http"
	"strings"
)

// BearerVerify Function check bearer token
func BearerVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		usr := GetUserFromRequest(r)

		if usr == nil {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode("Wrong auth token")
			err2.DebugErr(err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserFromRequest Fetch user model from request
func GetUserFromRequest(r *http.Request) *user.User {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		return nil
	}

	reqToken = splitToken[1]

	var usr *user.User
	db.MetaDb.GetConnection().Find(&usr, "token = ? ", reqToken)

	return usr
}
