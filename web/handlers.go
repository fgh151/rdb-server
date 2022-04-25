package web

import (
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/events"
	"db-server/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
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

func checkAccess(w http.ResponseWriter, r *http.Request) bool {
	topic := getTopic(r)
	p := models.Project{}.GetByTopic(topic).(models.Project)

	if !validateOrigin(p, r.Header.Get("Origin")) {
		send403Error(w, "Cors error. Origin not allowed")
		return false
	}

	if !validateKey(p.Key, r.Header.Get("db-key")) {
		send403Error(w, "db-key not Valid")
		return false
	}

	return true
}

func validateOrigin(p models.Project, origin string) bool {
	pOrigins := strings.Split(p.Origins, ";")
	for _, pOrigin := range pOrigins {
		if pOrigin == origin {
			return true
		}
	}

	return false
}

func validateKey(k1 string, k2 string) bool {
	return k1 == k2
}

func send403Error(w http.ResponseWriter, message string) {
	payload := map[string]string{"code": "not acceptable", "message": message}
	sendResponse(w, 403, payload, nil)
}

func PushHandler(w http.ResponseWriter, r *http.Request) {

	topic := getTopic(r)

	if checkAccess(w, r) {
		requestPayload := getPayload(r)
		_, err := drivers.GetDbInstance().Insert(os.Getenv("DB_NAME"), topic, requestPayload)

		var i interface{}
		sendResponse(w, 202, i, err)

		if err == nil {
			events.GetInstance().RegisterNewMessage(topic, requestPayload)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {

	topic := getTopic(r)

	vars := mux.Vars(r)
	rkey := vars["key"]

	if !validateKey(models.Project{}.GetKey(topic), rkey) {
		send403Error(w, "db-key not Valid")
	} else {
		c, err := upgrader.Upgrade(w, r, nil)

		events.GetInstance().Subscribe(topic, c)
		defer events.GetInstance().Unsubscribe(topic, c)

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
}

func FindHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)
	requestPayload := getPayload(r)

	if checkAccess(w, r) {

		res, err := drivers.GetDbInstance().Find(os.Getenv("DB_NAME"), topic, requestPayload)

		sendResponse(w, 200, res, err)
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)

	if checkAccess(w, r) {

		res, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic)

		sendResponse(w, 200, res, err)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	topic := getTopic(r)

	if checkAccess(w, r) {

		requestPayload := getPayload(r)

		vars := mux.Vars(r)
		id := vars["id"]

		res, err := drivers.GetDbInstance().Update(os.Getenv("DB_NAME"), topic, id, requestPayload)

		sendResponse(w, 202, res, err)
	}
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
