package project

import (
	"db-server/server/db"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (p Project) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var projects []Project

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &projects, make([]string, 0))

	y := make([]interface{}, len(projects))
	for i, v := range projects {
		y[i] = v
	}

	return y, nil
}

func (p Project) Total() *int64 {
	return db.MetaDb.TotalRecords(&Project{})
}

func (p Project) GetById(id string) (interface{}, error) {
	var project Project

	conn := db.MetaDb.GetConnection()

	tx := conn.First(&project, "id = ?", id)

	if tx.RowsAffected < 1 {
		return project, errors.New("no found")
	}

	return project, nil
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
