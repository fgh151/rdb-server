package models

import (
	"db-server/server"
	"github.com/google/uuid"
)

type Config struct {
	Id        uuid.UUID `gorm:"column:id;primary_key" json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	ProjectId uuid.UUID `json:"project_id"`
	Project   Project   `json:"project"`
}

func (p Config) List(limit int, offset int, sort string, order string) []interface{} {
	var configs []Config

	conn := server.MetaDb.GetConnection()

	conn.Find(&configs).Limit(limit).Offset(offset).Order(order + " " + sort)

	y := make([]interface{}, len(configs))
	for i, v := range configs {
		y[i] = v
	}

	return y
}

func (p Config) Total() *int64 {
	conn := server.MetaDb.GetConnection()
	var configs []Config
	var cnt int64
	conn.Find(&configs).Count(&cnt)

	return &cnt
}

func (p Config) GetById(id string) interface{} {
	var config Config

	conn := server.MetaDb.GetConnection()

	conn.Preload("Project").First(&config, "id = ?", id)

	return config
}

func (p Config) Delete(id string) {
	conn := server.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}
