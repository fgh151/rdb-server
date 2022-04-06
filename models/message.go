package models

import (
	"db-server/db"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	Id        uint `gorm:"primarykey" json:"id"`
	TopicId   uint
	Content   string
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t Message) GetById(id string) interface{} {
	var message Message

	conn := db.DB.GetConnection()

	conn.First(&message, "id = ?", id)

	return message
}

func (t Message) Find(filter map[string]interface{}) []Message {
	var messages []Message
	conn := db.DB.GetConnection()
	conn.Where(filter).Find(messages)
	return messages
}
