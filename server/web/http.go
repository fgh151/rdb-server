package web

import (
	"db-server/auth"
	"db-server/docs"
	"db-server/modules/push"
	"db-server/web"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

func StartServer(enableDocs *bool) {

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

	em.HandleFunc("/find/{topic}", web.FindHandler).Methods(http.MethodPost, http.MethodOptions)                // each request calls PushHandler
	em.HandleFunc("/list/{topic}", web.ListHandler).Methods(http.MethodGet, http.MethodOptions)                 // each request calls PushHandler
	em.HandleFunc("/subscribe/{topic}/{key}", web.SubscribeHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	em.HandleFunc("/{topic}", web.PushHandler).Methods(http.MethodPost, http.MethodOptions)                     // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", web.UpdateHandler).Methods(http.MethodPatch, http.MethodOptions)             // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", web.DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)            // each request calls PushHandler

	r.HandleFunc("/config/{id}", web.ApiConfigItem).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	r.HandleFunc("/dse/{id}", web.DSEItem).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	r.HandleFunc("/admin/auth", web.Auth).Methods(http.MethodPost, http.MethodOptions)          // each request calls PushHandler

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(auth.AdminVerify)
	admin.HandleFunc("/topics", web.ListTopics).Methods(http.MethodGet, http.MethodOptions)             // each request calls PushHandler
	admin.HandleFunc("/topics", web.CreateTopic).Methods(http.MethodPost, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/topics/{topic}/data", web.TopicData).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", web.TopicItem).Methods(http.MethodGet, http.MethodOptions)         // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", web.DeleteTopic).Methods(http.MethodDelete, http.MethodOptions)    // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", web.UpdateTopic).Methods(http.MethodPut, http.MethodOptions)       // each request calls PushHandler

	admin.HandleFunc("/users", web.ListUsers).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/users", web.CreateUser).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/users/{id}", web.UserItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/users/{id}", web.DeleteUser).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/users/{id}", web.UpdateUser).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/config", web.ListConfig).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/config", web.CreateConfig).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/config/{id}", web.ConfigItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/config/{id}", web.DeleteConfig).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/config/{id}", web.UpdateConfig).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/ds/dse/{dsId}", web.ListDse).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}", web.CreateDse).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", web.DseItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", web.DeleteDse).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/dse/{dsId}/{id}", web.UpdateDse).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/ds", web.ListDs).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/ds", web.CreateDs).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", web.DsItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", web.DeleteDs).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/ds/{id}", web.UpdateDs).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/cf", web.ListCf).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/cf", web.CreateCf).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", web.CfItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/cf/{id}/log", web.CfLog).Methods(http.MethodGet, http.MethodOptions)   // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", web.DeleteCf).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/cf/{id}", web.UpdateCf).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/pl", web.ListPipeline).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/pl", web.CreatePipeline).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", web.PipelineItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", web.DeletePipeline).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/pl/{id}", web.UpdatePipeline).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	push.AddAdminRoutes(admin)

	admin.HandleFunc("/cron", web.ListCron).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/cron", web.CreateCron).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", web.CronItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", web.DeleteCron).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/cron/{id}", web.UpdateCron).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler

	admin.HandleFunc("/em/list/{topic}", web.AdminListHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	r.HandleFunc("/api/user/oauth/{provider}/link", web.ApiOAuthLink).Methods(http.MethodGet, http.MethodOptions)   // each request calls PushHandler
	r.HandleFunc("/api/user/oauth/{provider}/{code}", web.ApiOAuthCode).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	r.HandleFunc("/api/user/auth", web.ApiAuth).Methods(http.MethodPost, http.MethodOptions)                             // each request calls PushHandler
	r.HandleFunc("/api/user/register", web.ApiRegister).Methods(http.MethodPost, http.MethodOptions)                     // each request calls PushHandler
	r.HandleFunc("/api/device/register", web.PushDeviceRegister).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	r.HandleFunc("/api/push/subscribe/{deviceId}", web.SubscribePushHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.BearerVerify)
	api.HandleFunc("/user/me", web.ApiMe).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	api.HandleFunc("/storage", web.StoragePut).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	api.HandleFunc("/cf/{id}/run", web.CfRun).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	api.HandleFunc("/cf/{id}/run/{rid}", web.CfRunLog).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

	push.AddApiRoutes(api)

	headersOk := handlers.AllowedHeaders(allowedHeaders)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions, http.MethodDelete})

	s := handlers.ExposedHeaders([]string{"X-Total-Count"})

	http.Handle("/", r)

	log.Debug("Start web server " + fullServer)
	log.Fatal(http.ListenAndServe(server, handlers.CORS(originsOk, headersOk, methodsOk, s)(r)))
}
