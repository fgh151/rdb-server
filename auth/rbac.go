package auth

import (
	"db-server/modules/user"
	"db-server/utils"
	"encoding/json"
	"net/http"
)

// AdminVerify Verify is request from admin user
func AdminVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		usr, err := utils.GetUserFromRequest(r)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode("Wrong auth token")
			return
		}

		if usr.(user.User).Admin != true {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode("Method not allowed")
			return
		}

		next.ServeHTTP(w, r)
	})
}
