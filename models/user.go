package models

import (
	"db-server/meta"
	"db-server/security"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id           uint           `gorm:"primarykey" json:"id"`
	Email        string         `gorm:"index" json:"email"`
	Token        string         `gorm:"index" json:"token"`
	PasswordHash string         `json:"-"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Model
}

func (p User) List() []interface{} {
	var users []User

	conn := meta.MetaDb.GetConnection()

	conn.Find(&users)

	y := make([]interface{}, len(users))
	for i, v := range users {
		y[i] = v
	}

	return y
}

func (p User) GetById(id string) interface{} {
	var user User

	conn := meta.MetaDb.GetConnection()

	conn.First(&user, "id = ?", id)

	return user
}

func (p User) Delete(id string) {
	conn := meta.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
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

	meta.MetaDb.GetConnection().Create(&u)

	return u
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
