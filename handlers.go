package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

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
