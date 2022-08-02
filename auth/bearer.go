package auth

import (
	err2 "db-server/err"
	"db-server/modules/user"
	"encoding/json"
	"net/http"
)

// BearerVerify Function check bearer token
func BearerVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := user.GetUserFromRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode("Wrong auth token e:" + err.Error())
			err2.DebugErr(err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
