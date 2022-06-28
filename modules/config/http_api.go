package config

import (
	err2 "db-server/err"
	"db-server/server/db"
	"db-server/utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AddAdminRoutes(admin *mux.Router) {
	admin.HandleFunc("/config", ListConfig).Methods(http.MethodGet, http.MethodOptions)           // each request calls PushHandler
	admin.HandleFunc("/config", CreateConfig).Methods(http.MethodPost, http.MethodOptions)        // each request calls PushHandler
	admin.HandleFunc("/config/{id}", ConfigItem).Methods(http.MethodGet, http.MethodOptions)      // each request calls PushHandler
	admin.HandleFunc("/config/{id}", DeleteConfig).Methods(http.MethodDelete, http.MethodOptions) // each request calls PushHandler
	admin.HandleFunc("/config/{id}", UpdateConfig).Methods(http.MethodPut, http.MethodOptions)    // each request calls PushHandler
}

func AddApiRoutes(api *mux.Router) {
	api.HandleFunc("/config/{id}", ApiConfigItem).Methods(http.MethodGet, http.MethodOptions) // each request calls PushHandler
}

// ListConfig godoc
// @Summary      List configs
// @Description  List configs
// @Tags         Config manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Security bearerAuth
// @Success      200  {array}   Config
//
// @Router       /admin/config [get]
func ListConfig(w http.ResponseWriter, r *http.Request) {
	utils.ListItems(Config{}, []string{}, r, w)
}

// ConfigItem godoc
// @Summary      Config info
// @Description  Config detail info
// @Tags         Config manager
// @tags Admin
// @Accept       json
// @Produce      json
// @Param        id path    string  true  "Config id" gg
// @Security bearerAuth
// @Success      200  {object}   Config
//
// @Router       /admin/config/{id} [get]
func ConfigItem(w http.ResponseWriter, r *http.Request) {
	utils.GetItem(Config{}, w, r)
}

// DeleteConfig godoc
// @Summary      Delete config
// @Description  Delete config
// @Tags         Config manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Config id" id
// @Security bearerAuth
// @Success      204
//
// @Router       /admin/config/{id} [delete]
func DeleteConfig(w http.ResponseWriter, r *http.Request) {
	utils.DeleteItem(Config{}, w, r)
}

// UpdateConfig
// @Summary      Update config
// @Description  Update config
// @Tags         Config manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        device    body     Config  true  "Config info" true
// @Param        id    path     string  true  "Config id" id
// @Success      200 {object} Config
// @Security bearerAuth
//
// @Router       /admin/config/{id} [put]
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	newm := Config{}

	err := json.NewDecoder(r.Body).Decode(&newm)
	db.MetaDb.GetConnection().Save(&newm)

	resp, _ := json.Marshal(newm)
	w.WriteHeader(200)
	_, err = w.Write(resp)
	err2.DebugErr(err)
}

// CreateConfig
// @Summary      Create config
// @Description  Create config
// @Tags         Config manager
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        config    body     user.CreateUserForm  true  "Config info" true
// @Success      200 {object} user.User
// @Security bearerAuth
//
// @Router       /admin/config [post]
func CreateConfig(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)
	model := Config{}

	err := json.NewDecoder(r.Body).Decode(&model)
	err2.DebugErr(err)
	model.Id, err = uuid.NewUUID()
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

// ApiConfigItem godoc
// @Summary      Config
// @Description  Get config by id
// @Tags         Config manager
// @Tags         Public Api
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "Config id" id
// @Param        db-key    header     string  true  "Auth key" gg
// @Success      200  {array}   interface{}
//
// @Router       /config/{id} [get]
func ApiConfigItem(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, r.RequestURI)

	vars := mux.Vars(r)
	model := Config{}.GetById(vars["id"]).(Config)

	rKey := r.Header.Get("db-key")

	if !utils.ValidateKey(model.Project.Key, rKey) {
		utils.Send403Error(w, "db-key not Valid")
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		_, err := w.Write([]byte(model.Body))
		err2.DebugErr(err)
	}
}
