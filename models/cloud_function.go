package models

import (
	"db-server/meta"
	"github.com/google/uuid"
)

type CloudFunction struct {
	Id        uuid.UUID `gorm:"primarykey" json:"id"`
	Title     string    `json:"title"`
	Container string    `json:"container"`
	Params    string    `json:"params"`
}

func (p CloudFunction) List(limit int, offset int, sort string, order string) []interface{} {
	var sources []CloudFunction

	conn := meta.MetaDb.GetConnection()

	conn.Find(&sources).Limit(limit).Offset(offset).Order(order + " " + sort)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func (p CloudFunction) Total() *int64 {
	conn := meta.MetaDb.GetConnection()
	var sources []CloudFunction
	var cnt int64
	conn.Find(&sources).Count(&cnt)

	return &cnt
}

func (p CloudFunction) GetById(id string) interface{} {
	var source CloudFunction

	conn := meta.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	return source
}

func (p CloudFunction) Delete(id string) {
	conn := meta.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}
