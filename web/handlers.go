package web

import (
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/events"
	"db-server/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetTopic(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["topic"]
}

func getPayload(r *http.Request) map[string]interface{} {
	var requestPayload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPayload)
	err2.PanicErr(err)

	return requestPayload
}

func checkAccess(w http.ResponseWriter, r *http.Request) bool {
	topic := GetTopic(r)
	p := models.Project{}.GetByTopic(topic).(models.Project)

	if !validateOrigin(p, r.Header.Get("Origin")) {
		Send403Error(w, "Cors error. Origin not allowed")
		return false
	}

	if !validateKey(p.Key, r.Header.Get("db-key")) {
		Send403Error(w, "db-key not Valid")
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

	log.Debug("Invalid origin")

	return false
}

func validateKey(k1 string, k2 string) bool {
	return k1 == k2
}

func Send403Error(w http.ResponseWriter, message string) {
	log.Debug("403 error")
	payload := map[string]string{"code": "not acceptable", "message": message}
	sendResponse(w, 403, payload, nil)
}

// PushHandler godoc
// @Summary      Create
// @Description  Create topic record
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Success      200  {array}   interface{}
//
// @Router       /em/{topic} [post]
func PushHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

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

// SubscribeHandler godoc
// @Summary      Subscribe
// @Description  Socket subscribe to topic
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Param        key    path     string  true  "Db key" string
// @Success      200  {array}   interface{}
//
// @Router       /em/subscribe/{topic}/{key} [get]
func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	vars := mux.Vars(r)
	rkey := vars["key"]

	if !validateKey(models.Project{}.GetKey(topic), rkey) {
		Send403Error(w, "db-key not Valid")
	} else {
		c, err := upgrader.Upgrade(w, r, nil)

		events.GetInstance().Subscribe(topic, c)
		defer events.GetInstance().Unsubscribe(topic, c)

		err = c.WriteMessage(1, []byte("test own message"))

		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer func() { _ = c.Close() }()
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

// SubscribePushHandler godoc
// @Summary      Subscribe
// @Description  Socket subscribe to push notifications
// @Tags         Push messages
// @Accept       json
// @Produce      json
// @Param        deviceId path    string  true  "Device id to subscribe" uuid
// @Success      200  {array}   interface{}
//
// @Router       /api/push/subscribe/{deviceId} [get]
func SubscribePushHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	c, err := upgrader.Upgrade(w, r, nil)

	events.GetPush().Subscribe(deviceId, c)
	defer events.GetPush().Unsubscribe(deviceId)

	err = c.WriteMessage(websocket.TextMessage, []byte("test own message"))

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() { _ = c.Close() }()
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

// FindHandler godoc
// @Summary      Search
// @Description  Search in topic
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Success      200  {array}   interface{}
//
// @Router       /em/find/{topic} [get]
func FindHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)
	requestPayload := getPayload(r)

	if checkAccess(w, r) {
		limit, offset, _, _ := GetPagination(r)

		res, err := drivers.GetDbInstance().Find(os.Getenv("DB_NAME"), topic, requestPayload, int64(limit), int64(offset))

		sendResponse(w, 200, res, err)
	}
}

// ListHandler godoc
// @Summary      List
// @Description  List topic records
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Success      200  {array}   interface{}
//
// @Router       /em/list/{topic} [get]
func ListHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	if checkAccess(w, r) {

		limit, offset, rorder, sort := GetPagination(r)

		v := r.URL.Query()
		filter := bson.D{{}}
		for _, param := range []string{"userId"} {
			if v.Has(param) {
				val := v.Get(param)
				if val != "" {
					filter = append(filter, primitive.E{Key: "userId", Value: val})
				}
			}
		}

		log.Debug("Mongo limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " order " + rorder + " sort " + sort)

		order, sort := drivers.GetMongoSort(sort, rorder)

		res, count, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic, int64(limit), int64(offset), order, sort, filter)

		w.Header().Add("X-Total-Count", strconv.FormatInt(count, 10))

		sendResponse(w, 200, res, err)
	}
}

// AdminListHandler godoc
// @Summary      List
// @Description  List topic records for admin access
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Security bearerAuth
// @Success      200  {array}   interface{}
//
// @Router       /admin/em/list/{topic} [get]
func AdminListHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	limit, offset, rorder, sort := GetPagination(r)

	v := r.URL.Query()
	filter := bson.D{{}}
	for _, param := range []string{"userId"} {
		if v.Has(param) {
			val := v.Get(param)
			if val != "" {
				filter = append(filter, primitive.E{Key: "userId", Value: val})
			}
		}
	}

	log.Debug("Mongo limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " order " + rorder + " sort " + sort)

	order, sort := drivers.GetMongoSort(sort, rorder)

	res, count, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic, int64(limit), int64(offset), order, sort, filter)

	w.Header().Add("X-Total-Count", strconv.FormatInt(count, 10))

	sendResponse(w, 200, res, err)
}

// UpdateHandler godoc
// @Summary      Update
// @Description  Update entity record
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Param        id    path     string  true  "Topic record id" id
// @Success      200  {array}   interface{}
//
// @Router       /em/{topic}/{id} [patch]
func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	if checkAccess(w, r) {

		requestPayload := getPayload(r)

		vars := mux.Vars(r)
		id := vars["id"]

		res, err := drivers.GetDbInstance().Update(os.Getenv("DB_NAME"), topic, id, requestPayload)

		sendResponse(w, 202, res, err)
	}
}

// DeleteHandler godoc
// @Summary      Delete
// @Description  Delete entity record
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" string
// @Param        id    path     string  true  "Topic record id" uuid
// @Success      200  {array}   interface{}
//
// @Router       /em/{topic}/{id} [delete]
func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	if checkAccess(w, r) {
		vars := mux.Vars(r)
		id := vars["id"]

		res, err := drivers.GetDbInstance().Delete(os.Getenv("DB_NAME"), topic, id)

		sendResponse(w, 202, res, err)
	}
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}, err error) {
	if err == nil {
		w.WriteHeader(statusCode)
		if payload != nil {
			resp, _ := json.Marshal(payload)
			_, err = w.Write(resp)
			err2.DebugErr(err)
		}
	} else {
		w.WriteHeader(500)
		_, err = w.Write([]byte(err.Error()))
		err2.DebugErr(err)
	}
}
