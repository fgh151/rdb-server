package em

import (
	"db-server/drivers"
	err2 "db-server/err"
	"db-server/events"
	"db-server/modules/project"
	"db-server/server"
	"db-server/server/db"
	"db-server/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func AddPublicApiRoutes(em *mux.Router) {
	em.HandleFunc("/find/{topic}", FindHandler).Methods(http.MethodPost, http.MethodOptions)                // each request calls PushHandler
	em.HandleFunc("/list/{topic}", ListHandler).Methods(http.MethodGet, http.MethodOptions)                 // each request calls PushHandler
	em.HandleFunc("/subscribe/{topic}/{key}", SubscribeHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	em.HandleFunc("/{topic}", PushHandler).Methods(http.MethodPost, http.MethodOptions)                     // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", UpdateHandler).Methods(http.MethodPatch, http.MethodOptions)             // each request calls PushHandler
	em.HandleFunc("/{topic}/{id}", DeleteHandler).Methods(http.MethodDelete, http.MethodOptions)            // each request calls PushHandler
}

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/topics", ListTopics).Methods(http.MethodGet, http.MethodOptions)             // each request calls PushHandler
	admin.HandleFunc("/topics", CreateTopic).Methods(http.MethodPost, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/topics/{topic}/data", TopicData).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", TopicItem).Methods(http.MethodGet, http.MethodOptions)         // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", DeleteTopic).Methods(http.MethodDelete, http.MethodOptions)    // each request calls PushHandler
	admin.HandleFunc("/topics/{id}", UpdateTopic).Methods(http.MethodPut, http.MethodOptions)       // each request calls PushHandler

	admin.HandleFunc("/em/list/{topic}", AdminListHandler).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
}

func GetTopic(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["topic"]
}

func checkAccess(w http.ResponseWriter, r *http.Request) bool {
	topic := GetTopic(r)
	p := project.Project{}.GetByTopic(topic).(project.Project)

	if !validateOrigin(p, r.Header.Get("Origin")) {
		utils.Send403Error(w, "Cors error. Origin not allowed")
		return false
	}

	if !utils.ValidateKey(p.Key, r.Header.Get("db-key")) {
		utils.Send403Error(w, "db-key not Valid")
		return false
	}

	return true
}

func validateOrigin(p project.Project, origin string) bool {
	pOrigins := strings.Split(p.Origins, ";")
	for _, pOrigin := range pOrigins {
		if pOrigin == origin {
			return true
		}
	}

	log.Debug("Invalid origin")

	return false
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
		requestPayload := utils.GetPayload(r)
		err := server.SaveTopicMessage(os.Getenv("DB_NAME"), topic, requestPayload)
		var i interface{}
		utils.SendResponse(w, 202, i, err)
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

	if !utils.ValidateKey(project.Project{}.GetKey(topic), rkey) {
		utils.Send403Error(w, "db-key not Valid")
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
	requestPayload := utils.GetPayload(r)

	if checkAccess(w, r) {
		limit, offset, _, _ := utils.GetPagination(r)

		res, err := drivers.GetDbInstance().Find(os.Getenv("DB_NAME"), topic, requestPayload, int64(limit), int64(offset))

		utils.SendResponse(w, 200, res, err)
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

		limit, offset, rorder, sort := utils.GetPagination(r)

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

		utils.SendResponse(w, 200, res, err)
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

	limit, offset, rorder, sort := utils.GetPagination(r)

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

	utils.SendResponse(w, 200, res, err)
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

		requestPayload := utils.GetPayload(r)

		vars := mux.Vars(r)
		id := vars["id"]

		res, err := drivers.GetDbInstance().Update(os.Getenv("DB_NAME"), topic, id, requestPayload)

		utils.SendResponse(w, 202, res, err)
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

		utils.SendResponse(w, 202, res, err)
	}
}

// ListTopics godoc
// @Summary      List topics
// @Description  List topics
// @Tags         TopicOutput
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   project.Project
//
// @Router       /admin/topics [get]
func ListTopics(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(project.Project{}, []string{}, r, w)
}

// TopicItem godoc
// @Summary      TopicOutput
// @Description  topic detail info
// @Tags         Entity manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "TopicOutput id" id
// @Security bearerAuth
// @Success      200  {object}   project.Project
//
// @Router       /admin/topics/{id} [get]
func TopicItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(project.Project{}, w, r)
}

// TopicData godoc
// @Summary      TopicOutput data
// @Description  topic data
// @Tags         Entity manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        topic path    string  true  "TopicOutput name"
// @Security bearerAuth
// @Success      200  {array} object
//
// @Router       /admin/topics/{topic}/data [get]
func TopicData(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	topic := GetTopic(r)

	limit, offset, rorder, sort := utils.GetPagination(r)

	order, sort := drivers.GetMongoSort(sort, rorder)

	log.Debug("Mongo limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " order " + rorder + " sort " + sort)

	res, count, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic, int64(limit), int64(offset), order, sort, bson.D{})

	var result []map[string]string

	for _, resArray := range res {
		record := make(map[string]string)
		for key, obj := range resArray.Map() {

			if key == "_id" {
				key = "id"
				obj = strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", obj), "bjectID(\"", ""), "\")", "")
			}
			record[key] = fmt.Sprintf("%v", obj)
		}
		result = append(result, record)
	}

	w.Header().Add("X-Total-Count", strconv.FormatInt(count, 10))

	utils.SendResponse(w, 200, result, err)
}

// UpdateTopic
// @Summary      Update topic
// @Description  Update topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     project.Project  true  "Project info" true
// @Param        id    path     string  true  "Project id" true
// @Success      200 {object} project.Project */
// @Security bearerAuth
//
// @Router       /admin/topics/{id} [put]
func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)

	var t = project.Project{}.GetById(vars["id"]).(project.Project)

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	db.MetaDb.GetConnection().Save(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// DeleteTopic godoc
// @Summary      Delete topic
// @Description  Delete topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "TopicOutput id" string
// @Success      204
//
// @Router       /admin/topics/{id} [delete]
func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(project.Project{}, w, r)
}

// CreateTopic
// @Summary      Create topic
// @Description  Create topic
// @Tags         Entity manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        topic    body     project.Project  true  "topic info" true
// @Success      200 {object} project.Project */
// @Security bearerAuth
//
// @Router       /admin/topics [post]
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	var t project.Project
	t.Id, _ = uuid.NewUUID()

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	db.MetaDb.GetConnection().Create(&t)

	resp, _ := json.Marshal(t)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}
