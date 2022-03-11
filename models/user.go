package models

import (
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

type CreateUserForm struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
