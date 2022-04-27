package models

import (
	"db-server/meta"
	"github.com/google/uuid"
)

type Config struct {
	Id    uuid.UUID `gorm:"column:id;primary_key" json:"id"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
}

func (c Config) List() []interface{} {
	var configs []Config

	conn := meta.MetaDb.GetConnection()

	conn.Find(&configs)

	y := make([]interface{}, len(configs))
	for i, v := range configs {
		y[i] = v
	}

	return y
}

func (p Config) GetById(id string) interface{} {
	var config Config

	conn := meta.MetaDb.GetConnection()

	conn.First(&config, "id = ?", id)

	return config
}

func (p Config) Delete(id string) {
	conn := meta.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}
