package oauth

import (
	"context"
	err2 "db-server/err"
	"db-server/models"
	"db-server/server"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/vk"
	"gorm.io/datatypes"
	"io"
	"strconv"
)

type UserOauth struct {
	// The record UUID
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`

	UserId    uuid.UUID
	User      models.User
	ServiceId string
	Data      datatypes.JSON
}

func (p UserOauth) GetUser() models.User {
	return models.User{}.GetById(p.UserId.String()).(models.User)
}

func (p UserOauth) GetByExternalId(id string) (UserOauth, error) {
	var user UserOauth

	conn := server.MetaDb.GetConnection()

	tx := conn.Preload("User").First(&user, "service_id = ?", id)

	if tx.RowsAffected > 0 {
		return user, nil
	}

	return user, errors.New("No user found")
}

type ClientOauth struct {
	Provider   string
	Config     *oauth2.Config
	GetUserApi string
}

func GetClient(provider string) (ClientOauth, error) {

	switch provider {
	case "github":
		return ClientOauth{
			Provider:   provider,
			GetUserApi: "https://api.github.com/user",
			Config: &oauth2.Config{
				ClientID:     models.GetAppSettingsByName("oauth_gh_client_id"),
				ClientSecret: models.GetAppSettingsByName("oauth_gh_client_secret"),
				RedirectURL:  models.GetAppSettingsByName("oauth_gh_client_redirect"),
				Scopes:       []string{"user"},
				Endpoint:     github.Endpoint,
			},
		}, nil
	case "vk":
		return ClientOauth{
			Provider:   provider,
			GetUserApi: "https://api.vk.com/method/users.get.json",
			Config: &oauth2.Config{
				ClientID:     models.GetAppSettingsByName("oauth_vk_client_id"),
				ClientSecret: models.GetAppSettingsByName("oauth_vk_client_secret"),
				RedirectURL:  models.GetAppSettingsByName("oauth_vk_client_redirect"),
				Scopes:       []string{"email"},
				Endpoint:     vk.Endpoint,
			},
		}, nil

	}

	return ClientOauth{}, errors.New("Provider " + provider + " not found")
}

func (c ClientOauth) GetUserByCode(code string) (UserOauth, error) {

	ctx := context.Background()

	tok, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return UserOauth{}, err
	}

	cl := c.Config.Client(ctx, tok)

	resp, err := cl.Get(c.GetUserApi)

	return GetOauthUser(resp.Body)
}

func GetOauthUser(resp io.ReadCloser) (UserOauth, error) {

	u := GithubUser{}
	decoder := json.NewDecoder(resp)
	err := decoder.Decode(&u)
	err2.DebugErr(err)

	bodyBytes, err := io.ReadAll(resp)
	err2.DebugErr(err)

	bodyString := string(bodyBytes)

	fmt.Println(bodyString)

	e, err := UserOauth{}.GetByExternalId(strconv.Itoa(u.Id))

	if err == nil {
		return e, nil
	}

	//u.Data = datatypes.JSON([]byte(bodyString))

	localUser, err := models.User{}.GetByEmail(u.Email)
	m, _ := json.Marshal(u)

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
