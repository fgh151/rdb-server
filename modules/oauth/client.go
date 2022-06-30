package oauth

import (
	"context"
	"db-server/modules/settings"
	"db-server/modules/user"
	"db-server/server/db"
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
	User      user.User
	ServiceId string
	Data      datatypes.JSON
}

// TableName Gorm table name
func (o UserOauth) TableName() string {
	return "user_oauth"
}

func (p UserOauth) GetUser() (user.User, error) {
	u, err := user.User{}.GetById(p.UserId.String())

	return u.(user.User), err
}

func (p UserOauth) GetByExternalId(id string) (UserOauth, error) {
	var u UserOauth

	conn := db.MetaDb.GetConnection()

	tx := conn.Preload("User").First(&u, "service_id = ?", id)

	if tx.RowsAffected > 0 {
		return u, nil
	}

	return u, errors.New("no found")
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

func GetClient(provider string, projectId uuid.UUID) (ClientOauth, error) {

	switch provider {
	case "github":
		return ClientOauth{
			Provider:     provider,
			GetUserApi:   "https://api.github.com/user",
			ExternalUser: GithubUser{},
			Config: &oauth2.Config{
				ClientID:     settings.GetAppSettingsByName(projectId, "oauth_gh_client_id"),
				ClientSecret: settings.GetAppSettingsByName(projectId, "oauth_gh_client_secret"),
				RedirectURL:  settings.GetAppSettingsByName(projectId, "oauth_gh_client_redirect"),
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
				ClientID:     settings.GetAppSettingsByName(projectId, "oauth_vk_client_id"),
				ClientSecret: settings.GetAppSettingsByName(projectId, "oauth_vk_client_secret"),
				RedirectURL:  settings.GetAppSettingsByName(projectId, "oauth_vk_client_redirect"),
				Scopes:       []string{"email"},
				Endpoint:     vk.Endpoint,
			},
		}, nil
	case "yandex":
		return ClientOauth{
			Provider:     provider,
			GetUserApi:   "https://login.yandex.ru/info?format=json",
			ExternalUser: YandexUser{},
			Config: &oauth2.Config{
				ClientID:     settings.GetAppSettingsByName(projectId, "oauth_yandex_client_id"),
				ClientSecret: settings.GetAppSettingsByName(projectId, "oauth_yandex_client_secret"),
				RedirectURL:  settings.GetAppSettingsByName(projectId, "oauth_yandex_client_redirect"),
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
