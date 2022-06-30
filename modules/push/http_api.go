package push

import (
	err2 "db-server/err"
	"db-server/events"
	"db-server/modules/push/device"
	"db-server/modules/push/models"
	"db-server/modules/user"
	"db-server/server/db"
	"db-server/utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/push", list).Methods(http.MethodGet, http.MethodOptions)               // each request calls PushHandler
	admin.HandleFunc("/push", create).Methods(http.MethodPost, http.MethodOptions)            // each request calls PushHandler
	admin.HandleFunc("/push/{id}", item).Methods(http.MethodGet, http.MethodOptions)          // each request calls PushHandler
	admin.HandleFunc("/push/{id}/run", item).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/push/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/push/{id}", update).Methods(http.MethodPut, http.MethodOptions)        // each request calls PushHandler
}

func AddApiRoutes(api *mux.Router) {
	api.HandleFunc("/push/{id}/run", run).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
}

func AddPublicApiRoutes(r *mux.Router) {
	r.HandleFunc("/api/push/subscribe/{deviceId}", subscribe).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler

}

// list godoc
// @Summary      List push messages
// @Description  List push messages
// @Tags         Push messages
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   models.PushMessage
//
// @Router       /admin/push [get]
func list(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(models.PushMessage{}, []string{}, r, w)
}

// create
// @Summary      Create push message
// @Description  Create push message
// @Tags         Push messages
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        push    body     models.PushMessage  true  "Push info" true
// @Success      200 {object} models.PushMessage
// @Security bearerAuth
//
// @Router       /admin/push [post]
func create(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := models.PushMessage{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	id, err := uuid.NewUUID()
	model.Id = id
	model.Sent = false

	err2.DebugErr(err)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 400)
		return
	}
	db.MetaDb.GetConnection().Create(&model)

	resp, _ := json.Marshal(model)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// item godoc
// @Summary      Push info
// @Description  Push detail info
// @Tags         Push messages
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Push id" id
// @Security bearerAuth
// @Success      200  {object}   models.PushMessage
//
// @Router       /admin/push/{id} [get]
func item(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(models.PushMessage{}, w, r)
}

// deleteItem godoc
// @Summary      Delete push
// @Description  Delete push
// @Tags         Push messages
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Push id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/push/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(models.PushMessage{}, w, r)
}

// update
// @Summary      Update push
// @Description  Update push
// @Tags         Push messages
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     models.PushMessage  true  "push info" true
// @Param        id    path     string  true  "push id" id
// @Success      200 {object} models.PushMessage
// @Security bearerAuth
//
// @Router       /admin/push/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	var exist = models.PushMessage{}.GetById(vars["id"]).(models.PushMessage)
	newm := models.PushMessage{}

	err := json.NewDecoder(r.Body).Decode(&newm)

	newm.CreatedAt = exist.CreatedAt

	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// run
// @Summary      Send push
// @Description  Send push with id
// @Tags         Push messages
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        db-key    header     string  true  "Auth key" true
// @Param        id    path     string  true  "Push id"
// @Success      200
//
// @Router       /api/push/{id}/run [get]
func run(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	vars := mux.Vars(r)
	message := models.PushMessage{}.GetById(vars["id"]).(models.PushMessage)
	go sendMessage(message)
	w.WriteHeader(200)
}

func sendMessage(p models.PushMessage) {

	for _, receiver := range p.Receivers {
		switch receiver.Device {
		case "ios":
			createPushLog(
				p,
				receiver,
				device.Ios{}.SendPush(p, receiver),
			)
			break

		case "android":
			createPushLog(
				p,
				receiver,
				device.Android{}.SendPush(p, receiver),
			)
			break
		default:
			createPushLog(p, receiver, device.InnerPush{}.SendPush(p, receiver))
		}
	}

	p.Sent = true
	p.SentAt = time.Now()
	db.MetaDb.GetConnection().Save(&p)
}

func createPushLog(message models.PushMessage, device user.UserDevice, err error) {

	id, _ := uuid.NewUUID()
	l := models.PushLog{
		Id:            id,
		PushMessageId: message.Id,
		UserDeviceId:  device.Id,
		Success:       err == nil,
		Error:         err.Error(),
		SentAt:        time.Now(),
	}

	db.MetaDb.GetConnection().Create(&l)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

// subscribe godoc
// @Summary      Subscribe
// @Description  Socket subscribe to push notifications
// @Tags         Push messages
// @Accept       json
// @Produce      json
// @Param        deviceId path    string  true  "Device id to subscribe" uuid
// @Success      200  {array}   interface{}
//
// @Router       /api/push/subscribe/{deviceId} [get]
func subscribe(w http.ResponseWriter, r *http.Request) {

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
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}
