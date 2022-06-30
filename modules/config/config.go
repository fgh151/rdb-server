package config

import (
	"db-server/modules/project"
	"db-server/server/db"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
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

func (p Config) List(limit int, offset int, sort string, order string, filter map[string]string) []interface{} {
	var configs []Config

	conn := db.MetaDb.GetConnection()

	conn.Limit(limit).Offset(offset).Order(clause.OrderBy{Expression: clause.Expr{SQL: "? ?", Vars: []interface{}{[]string{sort, order}}}}).Where(filter).Find(&configs)

	y := make([]interface{}, len(configs))
	for i, v := range configs {
		y[i] = v
	}

	return y
}

func (p Config) Total() *int64 {
	return db.MetaDb.TotalRecords(&Config{})
}

func (p Config) GetById(id string) interface{} {
	var config Config

	conn := db.MetaDb.GetConnection()

	conn.Preload("Project").First(&config, "id = ?", id)

	return config
}

func (p Config) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}
