package em

import (
	"db-server/drivers"
	"db-server/events"
	"db-server/modules/project"
	"db-server/modules/rdb"
	"db-server/server"
	"db-server/utils"
	"fmt"
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
	em.HandleFunc("/find/{topic}", find).Methods(http.MethodPost, http.MethodOptions)                // each request calls push
	em.HandleFunc("/list/{topic}", list).Methods(http.MethodGet, http.MethodOptions)                 // each request calls push
	em.HandleFunc("/subscribe/{topic}/{key}", subscribe).Methods(http.MethodGet, http.MethodOptions) // each request calls push
	em.HandleFunc("/{topic}", push).Methods(http.MethodPost, http.MethodOptions)                     // each request calls push
	em.HandleFunc("/{topic}/{id}", update).Methods(http.MethodPatch, http.MethodOptions)             // each request calls push
	em.HandleFunc("/{topic}/{id}", deleteItem).Methods(http.MethodDelete, http.MethodOptions)        // each request calls push
}

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/topics/{topic}/data", topicData).Methods(http.MethodGet, http.MethodOptions) // each request calls push
	admin.HandleFunc("/em/list/{topic}", adminList).Methods(http.MethodGet, http.MethodOptions)     // each request calls push
}

func getTopic(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["topic"]
}

func checkAccess(w http.ResponseWriter, r *http.Request) bool {
	topic := getTopic(r)

	dbi := rdb.Rdb{}.GetByCollection(topic)

	p := dbi.Project

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
		if pOrigin == "*" || pOrigin == origin {
			return true
		}
	}

	log.Debug("Invalid origin")

	return false
}

// push godoc
// @Summary      Create
// @Description  Create topic record
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Success      200  {array}   interface{}
//
// @Router       /em/{topic} [post]
func push(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

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

// subscribe godoc
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
func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

	vars := mux.Vars(r)
	rkey := vars["key"]

	dbi := rdb.Rdb{}.GetByCollection(topic)

	p := dbi.Project

	if !utils.ValidateKey(p.Key, rkey) {
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
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

// find godoc
// @Summary      Search
// @Description  Search in topic
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Success      200  {array}   interface{}
//
// @Router       /em/find/{topic} [get]
func find(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)
	requestPayload := utils.GetPayload(r)

	if checkAccess(w, r) {
		limit, offset, _, _ := utils.GetPagination(r)

		res, err := drivers.GetDbInstance().Find(os.Getenv("DB_NAME"), topic, requestPayload, int64(limit), int64(offset))

		utils.SendResponse(w, 200, res, err)
	}
}

// list godoc
// @Summary      List
// @Description  List topic records
// @Tags         Entity manager
// @Accept       json
// @Produce      json
// @Param        topic    path     string  true  "Topic name" gg
// @Success      200  {array}   interface{}
//
// @Router       /em/list/{topic} [get]

// sanitizeUserID trims and validates a userId value taken from the URL query
// before it is used in a MongoDB filter. It returns the sanitized value and
// a boolean indicating whether the value is acceptable.
func sanitizeUserID(raw string) (string, bool) {
	// Trim surrounding whitespace
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", false
	}
	// Optionally enforce a maximum length to avoid abuse
	if len(s) > 256 {
		return "", false
	}
	return s, true
}

func list(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

	if checkAccess(w, r) {

		limit, offset, rorder, sort := utils.GetPagination(r)

		v := r.URL.Query()
		filter := bson.D{{}}
		for _, param := range []string{"userId"} {
			if v.Has(param) {
				val := v.Get(param)
				if val != "" {
					if sanitized, ok := sanitizeUserID(val); ok {
						filter = append(filter, primitive.E{Key: "userId", Value: sanitized})
					}
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

// adminList godoc
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
func adminList(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

	limit, offset, rorder, sort := utils.GetPagination(r)

	v := r.URL.Query()
	filter := bson.D{{}}
	for _, param := range []string{"userId"} {
		if v.Has(param) {
			val := v.Get(param)
			if val != "" {
				if sanitized, ok := sanitizeUserID(val); ok {
					filter = append(filter, primitive.E{Key: "userId", Value: sanitized})
				}
			}
		}
	}

	log.Debug("Mongo limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " order " + rorder + " sort " + sort)

	order, sort := drivers.GetMongoSort(sort, rorder)

	res, count, err := drivers.GetDbInstance().List(os.Getenv("DB_NAME"), topic, int64(limit), int64(offset), order, sort, filter)

	w.Header().Add("X-Total-Count", strconv.FormatInt(count, 10))

	utils.SendResponse(w, 200, res, err)
}

// update godoc
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
func update(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

	if checkAccess(w, r) {

		requestPayload := utils.GetPayload(r)

		vars := mux.Vars(r)
		id := vars["id"]

		res, err := drivers.GetDbInstance().Update(os.Getenv("DB_NAME"), topic, id, requestPayload)

		utils.SendResponse(w, 202, res, err)
	}
}

// deleteItem godoc
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
func deleteItem(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

	if checkAccess(w, r) {
		vars := mux.Vars(r)
		id := vars["id"]

		res, err := drivers.GetDbInstance().Delete(os.Getenv("DB_NAME"), topic, id)

		utils.SendResponse(w, 202, res, err)
	}
}

// topicData godoc
// @Summary      Topic output data
// @Description  topic data
// @Tags         Entity manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        topic path    string  true  "Topic name"
// @Security bearerAuth
// @Success      200  {array} object
//
// @Router       /admin/topics/{topic}/data [get]
func topicData(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	topic := getTopic(r)

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
