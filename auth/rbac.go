package auth

import (
	"db-server/modules/user"
	"encoding/json"
	"net/http"
)

// AdminVerify Verify is request from admin user
func AdminVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		usr, err := user.GetUserFromRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode("Wrong auth token r: " + err.Error())
			return
		}

		if usr.Admin != true {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode("Method not allowed")
			return
		}

		next.ServeHTTP(w, r)
	})
}
