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

	godotenv.Load()

	client := GetDbInstance().getConnection()

	fmt.Println(client)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	init()

}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/push/{topic}", pushHandler) // each request calls pushHandler

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR")+":"+os.Getenv("SERVER_PORT"), nil))
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	_, err := GetDbInstance().insert(os.Getenv("DB_NAME"), topic, r.Body)

	if err == nil {
		w.WriteHeader(202)
	} else {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	fmt.Fprintf(w, "request %q\n", r.Body)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
