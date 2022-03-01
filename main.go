package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load()

	checkErr(err)

	client := GetDbInstance().getConnection()

	fmt.Println(client)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	initServer()
}

func initServer() {
	r := mux.NewRouter()
	r.HandleFunc("/push/{topic}", pushHandler) // each request calls pushHandler

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), nil))
}
