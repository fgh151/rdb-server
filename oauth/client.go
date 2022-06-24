package oauth

import (
	"context"
	"db-server/models"
	"db-server/server"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/vk"
	"golang.org/x/oauth2/yandex"
	"gorm.io/datatypes"
	"io"
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
	Provider     string
	Config       *oauth2.Config
	GetUserApi   string
	ExternalUser ExternalUser
}

type ExternalUser interface {
	GetOauthUser(resp io.ReadCloser) (UserOauth, error)
}

func GetClient(provider string) (ClientOauth, error) {

	switch provider {
	case "github":
		return ClientOauth{
			Provider:     provider,
			GetUserApi:   "https://api.github.com/user",
			ExternalUser: GithubUser{},
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
			Provider:     provider,
			GetUserApi:   "https://api.vk.com/method/users.get.json",
			ExternalUser: VkUser{},
			Config: &oauth2.Config{
				ClientID:     models.GetAppSettingsByName("oauth_vk_client_id"),
				ClientSecret: models.GetAppSettingsByName("oauth_vk_client_secret"),
				RedirectURL:  models.GetAppSettingsByName("oauth_vk_client_redirect"),
				Scopes:       []string{"email"},
				Endpoint:     vk.Endpoint,
			},
		}, nil
	case "yandex":
		return ClientOauth{
			Provider:     provider,
			GetUserApi:   "https://login.yandex.ru/info?format=json",
			ExternalUser: VkUser{},
			Config: &oauth2.Config{
				ClientID:     models.GetAppSettingsByName("oauth_yandex_client_id"),
				ClientSecret: models.GetAppSettingsByName("oauth_yandex_client_secret"),
				RedirectURL:  models.GetAppSettingsByName("oauth_yandex_client_redirect"),
				Scopes:       []string{"email", "profile"},
				Endpoint:     yandex.Endpoint,
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

	return c.ExternalUser.GetOauthUser(resp.Body)
}
