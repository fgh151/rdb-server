package web

import (
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/events"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

func getTopic(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["topic"]
}

func getPayload(r *http.Request) map[string]interface{} {
	var requestPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPayload)
	err2.CheckErr(err)

	return requestPayload
}

func PushHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)
	requestPayload := getPayload(r)

	_, err := drivers.GetDbInstance().Insert(os.Getenv("DB_NAME"), topic, requestPayload)

	var i interface{}
	sendResponse(w, 202, i, err)

	if err == nil {
		events.RegisterNewMessage(topic, requestPayload)
	}
}

var upgrader = websocket.Upgrader{} // use default options

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	events.Subscribe(getTopic(r), c)

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

func FindHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)
	requestPayload := getPayload(r)

	res, err := drivers.GetDbInstance().Find(os.Getenv("DB_NAME"), topic, requestPayload)

	sendResponse(w, 200, res, err)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)

	res, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic)

	sendResponse(w, 200, res, err)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)
	requestPayload := getPayload(r)

	id := requestPayload["id"]
	delete(requestPayload, "id")

	res, err := drivers.GetDbInstance().Update(os.Getenv("DB_NAME"), topic, id, requestPayload)

	sendResponse(w, 202, res, err)
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}, err error) {
	if err == nil {
		w.WriteHeader(statusCode)
		if payload != nil {
			resp, _ := json.Marshal(payload)
			w.Write(resp)
		}
	} else {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
