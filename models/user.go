package models

import (
	"db-server/db"
	"db-server/security"
	"errors"
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
}

func (p User) List() []interface{} {
	var users []User

	conn := db.DB.GetConnection()

	conn.Find(&users)

	y := make([]interface{}, len(users))
	for i, v := range users {
		y[i] = v
	}

	return y
}

func (p User) GetById(id string) interface{} {
	var user User

	conn := db.DB.GetConnection()

	conn.First(&user, "id = ?", id)

	return user
}

func (p User) Delete(id string) {
	conn := db.DB.GetConnection()
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

	db.DB.GetConnection().Create(&u)

	return u
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (f LoginForm) Login() (User, error) {
	var user User

	db.DB.GetConnection().First(&user, "email = ?", f.Email)

	if !user.ValidatePassword(f.Password) {
		return user, errors.New("invalid login or password")
	}

	return user, nil
}
