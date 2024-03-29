package config

import (
	"db-server/modules/project"
	"db-server/server/db"
	"errors"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Config struct {
	Id        uuid.UUID       `gorm:"column:id;primary_key" json:"id"`
	Title     string          `json:"title"`
	Body      datatypes.JSON  `json:"body"`
	ProjectId uuid.UUID       `json:"project_id"`
	Project   project.Project `json:"project"`
}

// TableName Gorm table name
func (p Config) TableName() string {
	return "config"
}

func (p Config) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var configs []Config

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &configs, make([]string, 0))

	y := make([]interface{}, len(configs))
	for i, v := range configs {
		y[i] = v
	}

	return y, nil
}

func (p Config) Total() *int64 {
	return db.MetaDb.TotalRecords(&Config{})
}

func (p Config) GetById(id string) (interface{}, error) {
	var config Config

	conn := db.MetaDb.GetConnection()

	tx := conn.Preload("Project").First(&config, "id = ?", id)

	if tx.RowsAffected < 1 {
		return config, errors.New("no found")
	}

	return config, nil
}

func (p Config) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}
