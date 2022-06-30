package oauth

import (
	"crypto/rand"
	err2 "db-server/err"
	"db-server/modules/project"
	"db-server/utils"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func AddPublicApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/user/oauth/{provider}/link", ApiOAuthLink).Methods(http.MethodGet, http.MethodOptions)   // each request calls PushHandler
	r.HandleFunc("/api/user/oauth/{provider}/{code}", ApiOAuthCode).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

}

// ApiOAuthLink godoc
// @Summary      OAuth link
// @Description  Get link for oauth
// @Tags         OAuth
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        provider    path     string  true  "Provider name"
// @Param        db-key    header     string  true  "Auth key" gg
// @Success      200  {string} string
//
// @Router       /api/user/oauth/{provider}/link [get]
func ApiOAuthLink(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	provider := vars["provider"]

	rKey := r.Header.Get("db-key")
	p, err := project.Project{}.GetByKey(rKey)

	if err != nil {
		payload := map[string]string{"code": "not acceptable", "message": err.Error()}
		utils.SendResponse(w, 500, payload, nil)
	}

	client, _ := GetClient(provider, p.Id)

	state := generateStateOauthCookie()
	url := client.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	w.WriteHeader(200)
	w.Write([]byte(url))
}

func generateStateOauthCookie() string {
	b := make([]byte, 128)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

// ApiOAuthCode godoc
// @Summary      OAuth user
// @Description  Get user by oauth code
// @Tags         OAuth
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        provider    path     string  true  "Provider name"
// @Param        code    path     string  true  "Code"
// @Param        db-key    header     string  true  "Auth key" gg
// @Success      200  {object} user.User
//
// @Router       /api/user/oauth/{provider}/{code} [get]
func ApiOAuthCode(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	provider := vars["provider"]
	code := vars["code"]

	rKey := r.Header.Get("db-key")
	p, err := project.Project{}.GetByKey(rKey)

	if err != nil {
		payload := map[string]string{"code": "not acceptable", "message": err.Error()}
		utils.SendResponse(w, 500, payload, nil)
	}

	client, _ := GetClient(provider, p.Id)

	u, err := client.GetUserByCode(code)

	if err != nil {
		log.Debug(err)
		payload := map[string]string{"code": "not acceptable", "message": err.Error()}
		utils.SendResponse(w, 500, payload, nil)
		return
	}

	usr := u.GetUser()
	usr.UpdateLastLogin()

	rresp, _ := json.Marshal(usr)
	w.WriteHeader(200)
	_, err = w.Write(rresp)
	err2.DebugErr(err)

	return
}
