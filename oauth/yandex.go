package oauth

import (
	err2 "db-server/err"
	"db-server/models"
	"db-server/server"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"strconv"
)

type YandexUser struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"real_name"`
	Email string `json:"default_email"`
}

func (u YandexUser) GetOauthUser(resp io.ReadCloser) (UserOauth, error) {

	data := make(map[string]interface{})

	decoder := json.NewDecoder(resp)
	err := decoder.Decode(&u)
	err2.DebugErr(err)
	err = decoder.Decode(&data)
	err2.DebugErr(err)

	e, err := UserOauth{}.GetByExternalId(strconv.Itoa(u.Id))

	if err == nil {
		return e, nil
	}

	localUser, err := models.User{}.GetByEmail(u.Email)
	m, _ := json.Marshal(data)

	id, _ := uuid.NewUUID()
	ru := UserOauth{
		Id:        id,
		Data:      m,
		UserId:    localUser.(models.User).Id,
		ServiceId: strconv.Itoa(u.Id),
	}

	server.MetaDb.GetConnection().Create(ru)

	return ru, err
}
