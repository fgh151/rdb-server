package web

import (
	"db-server/auth"
	"db-server/docs"
	"db-server/modules/cf"
	"db-server/modules/config"
	"db-server/modules/cron"
	"db-server/modules/ds"
	"db-server/modules/em"
	"db-server/modules/oauth"
	"db-server/modules/pipeline"
	"db-server/modules/project"
	"db-server/modules/push"
	"db-server/modules/storage"
	"db-server/modules/user"
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

	emr := r.PathPrefix("/em").Subrouter()

	em.AddPublicApiRoutes(emr)
	config.AddApiRoutes(r)
	ds.AddPublicApiRoutes(r)
	user.AddPublicApiRoutes(r)

	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(auth.AdminVerify)

	project.AddAdminRoutes(admin)
	em.AddAdminRoutes(admin)
	user.AddAdminRoutes(admin)
	config.AddAdminRoutes(admin)
	ds.AddAdminRoutes(admin)
	cf.AddAdminRoutes(admin)
	pipeline.AddAdminRoutes(admin)
	push.AddAdminRoutes(admin)
	cron.AddAdminRoutes(admin)

	push.AddPublicApiRoutes(r)
	oauth.AddPublicApiRoutes(r)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.BearerVerify)
	user.AddApiRoutes(api)
	storage.AddApiRoutes(api)
	push.AddApiRoutes(api)
	cf.AddApiRoutes(api)

	headersOk := handlers.AllowedHeaders(allowedHeaders)
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions, http.MethodDelete})

	s := handlers.ExposedHeaders([]string{"X-Total-Count"})

	http.Handle("/", r)

	log.Debug("Start web server " + fullServer)
	log.Fatal(http.ListenAndServe(server, handlers.CORS(originsOk, headersOk, methodsOk, s)(r)))
}
