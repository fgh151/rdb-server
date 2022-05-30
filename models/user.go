package models

import (
	"db-server/security"
	"db-server/server"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id           uuid.UUID      `gorm:"primarykey" json:"id"`
	Email        string         `gorm:"index" json:"email"`
	Token        string         `gorm:"index" json:"token"`
	PasswordHash string         `json:"-"`
	Admin        bool           `gorm:"index;default:false;type:bool" json:"admin"`
	Active       bool           `gorm:"index;default:true;type:bool" json:"active"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	LastLogin    *time.Time     `json:"lastLogin,omitempty"`
}

func (p User) List(limit int, offset int, sort string, order string) []interface{} {
	var users []User

	conn := server.MetaDb.GetConnection()

	conn.Find(&users).Limit(limit).Offset(offset).Order(order + " " + sort)

	y := make([]interface{}, len(users))
	for i, v := range users {
		y[i] = v
	}

	return y
}

func (p User) Total() *int64 {
	conn := server.MetaDb.GetConnection()
	var users []User
	var cnt int64
	conn.Find(&users).Count(&cnt)

	return &cnt
}

func (p User) GetById(id string) interface{} {
	var user User

	conn := server.MetaDb.GetConnection()

	conn.First(&user, "id = ?", id)

	return user
}

func (p User) Delete(id string) {
	conn := server.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p User) ValidatePassword(password string) bool {
	return p.PasswordHash == security.HashPassword(password)
}

type CreateUserForm struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

func (f CreateUserForm) Save() User {
	var u = User{
		Email:        f.Email,
		PasswordHash: security.HashPassword(f.Password),
		Token:        security.GenerateRandomString(15),
	}

	u.Id, _ = uuid.NewUUID()

	server.MetaDb.GetConnection().Create(&u)

	return u
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (f LoginForm) AdminLogin() (User, error) {
	return f.login(&User{Email: f.Email, Admin: true, Active: true})
}

func (f LoginForm) ApiLogin() (User, error) {
	return f.login(&User{Email: f.Email, Active: true})
}

func (f LoginForm) login(condition *User) (User, error) {
	var login User

	server.MetaDb.GetConnection().Where(condition).First(&login)

	if !login.ValidatePassword(f.Password) {
		return login, errors.New("invalid login or password")
	}

	now := time.Now()

	login.LastLogin = &now
	server.MetaDb.GetConnection().Save(&login)

	return login, nil
}
