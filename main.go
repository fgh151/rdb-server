package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	err := godotenv.Load()
	checkErr(err)

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	client := GetDbInstance().getConnection()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//socket.StartSocketServer()

	initServer()
}

func initServer() {
	r := mux.NewRouter()
	r.HandleFunc("/push/{topic}", pushHandler)           // each request calls pushHandler
	r.HandleFunc("/subscribe/{topic}", subscribeHandler) // each request calls pushHandler

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), nil))
}
