package main

import (
	"db-server/events"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func pushHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	var requestPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPayload)
	checkErr(err)

	_, err = GetDbInstance().insert(os.Getenv("DB_NAME"), topic, requestPayload)

	if err == nil {
		w.WriteHeader(202)
		events.RegisterNewMessage(topic, requestPayload)
	} else {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	fmt.Fprintf(w, "request %q\n", r.Body)
}
