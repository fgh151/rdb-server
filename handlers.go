package main

import (
	"db-server/events"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
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

var upgrader = websocket.Upgrader{} // use default options

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	c, err := upgrader.Upgrade(w, r, nil)

	events.Subscribe(topic, c)

	c.WriteMessage(1, []byte("test own message"))

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
