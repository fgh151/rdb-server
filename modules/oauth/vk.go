package oauth

import (
	err2 "db-server/err"
	"db-server/modules/user"
	"db-server/server/db"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"strconv"
)

type VkUser struct {
	Login string `json:"login"`
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (u VkUser) GetOauthUser(resp io.ReadCloser) (UserOauth, error) {

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

	//u.Data = datatypes.JSON([]byte(bodyString))

	localUser, err := user.User{}.GetByEmail(u.Email)
	m, _ := json.Marshal(data)

	id, _ := uuid.NewUUID()
	ru := UserOauth{
		Id:        id,
		Data:      m,
		UserId:    localUser.(user.User).Id,
		ServiceId: strconv.Itoa(u.Id),
	}

	db.MetaDb.GetConnection().Create(ru)

	return ru, err
}
