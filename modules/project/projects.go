package project

import (
	"db-server/server/db"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// swagger:model
type Project struct {
	Id        uuid.UUID      `gorm:"primarykey" json:"id"`
	Name      string         `json:"name"`
	Key       string         `json:"key"`
	Origins   string         `json:"origins"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p Project) List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{} {
	var projects []Project

	conn := db.MetaDb.GetConnection()

	conn.Offset(offset).Limit(limit).Order(clause.OrderBy{Expression: clause.Expr{SQL: "? ?", Vars: []interface{}{[]string{sort, order}}}}).Where(filter).Find(&projects)

	y := make([]interface{}, len(projects))
	for i, v := range projects {
		y[i] = v
	}

	return y
}

func (p Project) Total() *int64 {
	return db.MetaDb.TotalRecords(&Project{})
}

func (p Project) GetById(id string) interface{} {
	var project Project

	conn := db.MetaDb.GetConnection()

	conn.First(&project, "id = ?", id)

	return project
}

func (p Project) GetByKey(key string) (Project, error) {
	var project Project

	conn := db.MetaDb.GetConnection()

	tx := conn.First(&project, "key = ?", key)

	if tx.RowsAffected < 1 {
		return project, errors.New("No project found")
	}

	return project, nil
}

func (p Project) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

// TableName Gorm table name
func (p Project) TableName() string {
	return "project"
}
