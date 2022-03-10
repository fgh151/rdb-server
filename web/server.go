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
	r.HandleFunc("/push/{topic}", PushHandler).Methods("POST")                // each request calls PushHandler
	r.HandleFunc("/push/{topic}", UpdateHandler).Methods("PATCH")             // each request calls PushHandler
	r.HandleFunc("/find/{topic}", FindHandler).Methods("POST")                // each request calls PushHandler
	r.HandleFunc("/list/{topic}", ListHandler).Methods("GET")                 // each request calls PushHandler
	r.HandleFunc("/subscribe/{topic}/{key}", SubscribeHandler).Methods("GET") // each request calls PushHandler

	r.HandleFunc("/admin/topics", ListTopics).Methods("GET")     // each request calls PushHandler
	r.HandleFunc("/admin/topics/{id}", TopicItem).Methods("GET") // each request calls PushHandler

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	s := handlers.ExposedHeaders([]string{"X-Total-Count"})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), handlers.CORS(originsOk, headersOk, methodsOk, s)(r)))
}
