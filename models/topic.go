package models

import (
	"db-server/db"
	"gorm.io/gorm"
	"time"
)

type Topic struct {
	Id        uint `gorm:"primarykey" json:"id"`
	Title     string
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Messages  []Message
}

func (t Topic) GetById(id string) interface{} {
	var project Topic

	conn := db.DB.GetConnection()

	conn.First(&project, "id = ?", id)

	return project
}

func (t Topic) GetByTitle(title string) interface{} {
	var topic Topic

	conn := db.DB.GetConnection()

	res := conn.First(&topic, "title = ?", title)

	if res.RowsAffected < 1 {
		topic.Title = title
		conn.Create(&topic)
	}

	return topic
}

func (t Topic) Delete(id string) {
	conn := db.DB.GetConnection()
	conn.Where("id = ?", id).Delete(&t)
}

func (t Topic) List() []interface{} {
	var topics []Topic

	conn := db.DB.GetConnection()

	conn.Find(&topics)

	y := make([]interface{}, len(topics))
	for i, v := range topics {
		y[i] = v
	}

	return y
}
