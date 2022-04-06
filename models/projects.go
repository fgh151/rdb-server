package models

import (
	"db-server/db"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	Id        uint           `gorm:"primarykey" json:"id"`
	Topic     string         `json:"topic"`
	Key       string         `json:"key"`
	Origins   string         `json:"origins"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p Project) List() []interface{} {
	var projects []Project

	conn := db.DB.GetConnection()

	conn.Find(&projects)

	y := make([]interface{}, len(projects))
	for i, v := range projects {
		y[i] = v
	}

	return y
}

func (p Project) GetById(id string) interface{} {
	var project Project

	conn := db.DB.GetConnection()

	conn.First(&project, "id = ?", id)

	return project
}

func (p Project) GetByTopic(topic string) interface{} {
	var project Project

	conn := db.DB.GetConnection()

	conn.First(&project, "topic = ?", topic)

	return project
}

func (p Project) Delete(id string) {
	conn := db.DB.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (c Project) GetKey(topic string) string {
	var project Project

	conn := db.DB.GetConnection()

	conn.First(&project, "topic = ?", topic)

	return project.Key
}
