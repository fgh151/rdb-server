package web

import (
	"github.com/gorilla/mux"
	"log"
	http2 "net/http"
	"os"
)

func InitServer() {
	r := mux.NewRouter()
	r.HandleFunc("/push/{topic}", PushHandler)           // each request calls PushHandler
	r.HandleFunc("/subscribe/{topic}", SubscribeHandler) // each request calls PushHandler

	http2.Handle("/", r)
	log.Fatal(http2.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), nil))
}
