package web

import (
	"db-server/auth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func InitServer() {

	allowedHeaders := []string{
		"Access-Control-Allow-Origin",
		"Content-Type",
		"X-Requested-With",
		"Accept",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Accept-Language",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Cache-Control",
		"Connection",
		"Host",
		"Origin",
		"Sec-Fetch-Dest",
		"Sec-Fetch-Mode",
		"Sec-Fetch-Site",
		"User-Ag",
		"db-key",
	}

	r := mux.NewRouter()

	em := r.PathPrefix("/em").Subrouter()

	em.HandleFunc("/find/{topic}", FindHandler).Methods(http.MethodPost, http.MethodOptions)     // each request calls PushHandler
	em.HandleFunc("/list/{topic}", ListHandler).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	em.HandleFunc("/subscribe/{topic}/{key}", SubscribeHandler).Methods(http.MethodGet)          // each request calls PushHandler
	em.HandleFunc("/{topic}", PushHandler).Methods(http.MethodPost, http.MethodOptions)          // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", UpdateHandler).Methods(http.MethodPatch, http.MethodOptions)  // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", DeleteHandler).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler

	r.HandleFunc("/config/{id}", ApiConfigItem).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	r.HandleFunc("/dse/{id}", DSEItem).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	r.HandleFunc("/admin/auth", Auth).Methods(http.MethodPost, http.MethodOptions)          // each request calls PushHandler

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(auth.AdminVerify)
	admin.HandleFunc("/topics", ListTopics).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/topics", CreateTopic).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", TopicItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", DeleteTopic).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", UpdateTopic).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/users", ListUsers).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/users", CreateUser).Methods(http.MethodPost, http.MethodOptions)         // each request calls PushHandler
	admin.HandleFunc("/users/{id}", UserItem).Methods(http.MethodGet, http.MethodOptions)       // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", DeleteUser).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", UpdateUser).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/config", ListConfig).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/config", CreateConfig).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/config/{id}", ConfigItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/config/{id}", DeleteConfig).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/config/{id}", UpdateConfig).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	//admin.HandleFunc("/ds/e", ListDs).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	//admin.HandleFunc("/ds/e", CreateDs).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	//admin.HandleFunc("/ds/{id}/e", DsItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	//admin.HandleFunc("/ds/{id}/e", DeleteDs).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	//admin.HandleFunc("/ds/{id}/e", UpdateDs).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/ds", ListDs).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds", CreateDs).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", DsItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", DeleteDs).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", UpdateDs).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	r.HandleFunc("/api/user/auth", ApiAuth).Methods(http.MethodPost, http.MethodOptions)         // each request calls PushHandler
	r.HandleFunc("/api/user/register", ApiRegister).Methods(http.MethodPost, http.MethodOptions) // each request calls PushHandler
	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.BearerVerify)
	api.HandleFunc("/user/me", ApiMe).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	api.HandleFunc("/storage", StoragePut).Methods(http.MethodPost, http.MethodOptions) // each request calls PushHandler

	headersOk := handlers.AllowedHeaders(allowedHeaders)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions, http.MethodDelete})

	s := handlers.ExposedHeaders([]string{"X-Total-Count"})

	http.Handle("/", r)

	log.Debug("Start web server")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), handlers.CORS(originsOk, headersOk, methodsOk, s)(r)))
}
