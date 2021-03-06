package web

import (
	"db-server/auth"
	"db-server/docs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

func InitServer(enableDocs *bool) {

	server := os.Getenv("SERVER_ADDR") + ":" + os.Getenv("SERVER_PORT")
	fullServer := os.Getenv("SERVER_SCHEME") + "://" + server

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

	if *enableDocs {
		docs.SwaggerInfo.Title = "Db server API"
		docs.SwaggerInfo.Description = "Db server API description."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fullServer
		docs.SwaggerInfo.BasePath = "/"
		docs.SwaggerInfo.Schemes = []string{os.Getenv("SERVER_SCHEME")}

		r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
			httpSwagger.URL(fullServer+"/swagger/doc.json"), //The url pointing to API definition
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)).Methods(http.MethodGet)
	}

	em := r.PathPrefix("/em").Subrouter()

	em.HandleFunc("/find/{topic}", FindHandler).Methods(http.MethodPost, http.MethodOptions)                // each request calls PushHandler
	em.HandleFunc("/list/{topic}", ListHandler).Methods(http.MethodGet, http.MethodOptions)                 // each request calls PushHandler
	em.HandleFunc("/subscribe/{topic}/{key}", SubscribeHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	em.HandleFunc("/{topic}", PushHandler).Methods(http.MethodPost, http.MethodOptions)                     // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", UpdateHandler).Methods(http.MethodPatch, http.MethodOptions)             // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)            // each request calls PushHandler

	r.HandleFunc("/config/{id}", ApiConfigItem).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	r.HandleFunc("/dse/{id}", DSEItem).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	r.HandleFunc("/admin/auth", Auth).Methods(http.MethodPost, http.MethodOptions)          // each request calls PushHandler

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(auth.AdminVerify)
	admin.HandleFunc("/topics", ListTopics).Methods(http.MethodGet, http.MethodOptions)             // each request calls PushHandler
	admin.HandleFunc("/topics", CreateTopic).Methods(http.MethodPost, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/topics/{topic}/data", TopicData).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", TopicItem).Methods(http.MethodGet, http.MethodOptions)         // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", DeleteTopic).Methods(http.MethodDelete, http.MethodOptions)    // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", UpdateTopic).Methods(http.MethodPut, http.MethodOptions)       // each request calls PushHandler

	admin.HandleFunc("/users", ListUsers).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/users", CreateUser).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/users/{id}", UserItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/users/{id}", DeleteUser).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/users/{id}", UpdateUser).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/config", ListConfig).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/config", CreateConfig).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/config/{id}", ConfigItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/config/{id}", DeleteConfig).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/config/{id}", UpdateConfig).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/ds/dse/{dsId}", ListDse).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}", CreateDse).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", DseItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", DeleteDse).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", UpdateDse).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/ds", ListDs).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds", CreateDs).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", DsItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", DeleteDs).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", UpdateDs).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/cf", ListCf).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/cf", CreateCf).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", CfItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/cf/{id}/log", CfLog).Methods(http.MethodGet, http.MethodOptions)   // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", DeleteCf).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", UpdateCf).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/pl", ListPipeline).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/pl", CreatePipeline).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", PipelineItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", DeletePipeline).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", UpdatePipeline).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/push", ListPush).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/push", CreatePush).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/push/{id}", PushItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/push/{id}/run", PushItem).Methods(http.MethodGet, http.MethodOptions)  // each request calls PushHandler
	admin.HandleFunc("/push/{id}", DeletePush).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/push/{id}", UpdatePush).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/cron", ListCron).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/cron", CreateCron).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", CronItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", DeleteCron).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", UpdateCron).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/settings/{projectId}/oauth", SettingsOauth).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	admin.HandleFunc("/em/list/{topic}", AdminListHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	r.HandleFunc("/api/user/oauth/{provider}/link", ApiOAuthLink).Methods(http.MethodGet, http.MethodOptions)   // each request calls PushHandler
	r.HandleFunc("/api/user/oauth/{provider}/{code}", ApiOAuthCode).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	r.HandleFunc("/api/user/auth", ApiAuth).Methods(http.MethodPost, http.MethodOptions)                             // each request calls PushHandler
	r.HandleFunc("/api/user/register", ApiRegister).Methods(http.MethodPost, http.MethodOptions)                     // each request calls PushHandler
	r.HandleFunc("/api/device/register", PushDeviceRegister).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	r.HandleFunc("/api/push/subscribe/{deviceId}", SubscribePushHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.BearerVerify)
	api.HandleFunc("/user/me", ApiMe).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	api.HandleFunc("/storage", StoragePut).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	api.HandleFunc("/cf/{id}/run", CfRun).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	api.HandleFunc("/cf/{id}/run/{rid}", CfRunLog).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	api.HandleFunc("/push/{id}/run", PushRun).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler

	headersOk := handlers.AllowedHeaders(allowedHeaders)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions, http.MethodDelete})

	s := handlers.ExposedHeaders([]string{"X-Total-Count"})

	http.Handle("/", r)

	log.Debug("Start web server " + fullServer)
	log.Fatal(http.ListenAndServe(server, handlers.CORS(originsOk, headersOk, methodsOk, s)(r)))
}
