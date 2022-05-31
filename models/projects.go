package models

import (
	"db-server/server"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Project struct {
	Id        uuid.UUID      `gorm:"primarykey" json:"id"`
	Topic     string         `json:"topic"`
	Key       string         `json:"key"`
	Origins   string         `json:"origins"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p Project) List(limit int, offset int, sort string, order string) []interface{} {
	var projects []Project

	conn := server.MetaDb.GetConnection()

	conn.Offset(offset).Limit(limit).Order(order + " " + sort).Find(&projects)

	y := make([]interface{}, len(projects))
	for i, v := range projects {
		y[i] = v
	}

	return y
}

func (p Project) Total() *int64 {
	return TotalRecords(&Project{})
}

func (p Project) GetById(id string) interface{} {
	var project Project

	conn := server.MetaDb.GetConnection()

	conn.First(&project, "id = ?", id)

	return project
}

func (p Project) GetByTopic(topic string) interface{} {
	var project Project

	conn := server.MetaDb.GetConnection()

	conn.First(&project, "topic = ?", topic)

	return project
}

func (p Project) Delete(id string) {
	conn := server.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (c Project) GetKey(topic string) string {
	var project Project

	conn := server.MetaDb.GetConnection()

	conn.First(&project, "topic = ?", topic)

	return project.Key
}
