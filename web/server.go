package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func InitServer() {
	r := mux.NewRouter()
	r.HandleFunc("/push/{topic}", PushHandler).Methods(http.MethodPost)                // each request calls PushHandler
	r.HandleFunc("/push/{topic}", UpdateHandler).Methods(http.MethodPatch)             // each request calls PushHandler
	r.HandleFunc("/find/{topic}", FindHandler).Methods(http.MethodPost)                // each request calls PushHandler
	r.HandleFunc("/list/{topic}", ListHandler).Methods(http.MethodGet)                 // each request calls PushHandler
	r.HandleFunc("/subscribe/{topic}/{key}", SubscribeHandler).Methods(http.MethodGet) // each request calls PushHandler

	r.HandleFunc("/admin/topics", ListTopics).Methods(http.MethodGet)                       // each request calls PushHandler
	r.HandleFunc("/admin/topics", CreateTopic).Methods(http.MethodPost, http.MethodOptions) // each request calls PushHandler
	r.HandleFunc("/admin/topics/{id}", TopicItem).Methods(http.MethodGet)                   // each request calls PushHandler
	r.HandleFunc("/admin/topics/{id}", DeleteTopic).Methods(http.MethodDelete)              // each request calls PushHandler
	//r.Use(mux.CORSMethodMiddleware(r))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPatch, http.MethodOptions, http.MethodDelete})

	s := handlers.ExposedHeaders([]string{"X-Total-Count"})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), handlers.CORS(originsOk, headersOk, methodsOk, s)(r)))
}
