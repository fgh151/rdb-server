package auth

import (
	"encoding/json"
	"net/http"
)

// AdminVerify Verify is request from admin user
func AdminVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := GetUserFromRequest(r)

		if user == nil {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode("Wrong auth token")
			return
		}

		if user.Admin != true {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode("Method not allowed")
			return
		}

		next.ServeHTTP(w, r)
	})
}
