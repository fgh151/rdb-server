package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func InitServer() {
	r := mux.NewRouter()
	r.HandleFunc("/push/{topic}", PushHandler).Methods("POST")          // each request calls PushHandler
	r.HandleFunc("/find/{topic}", FindHandler).Methods("POST")          // each request calls PushHandler
	r.HandleFunc("/subscribe/{topic}", SubscribeHandler).Methods("GET") // each request calls PushHandler

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), nil))
}
