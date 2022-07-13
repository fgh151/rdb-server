package rdb

import (
	"db-server/modules/project"
	"db-server/server/db"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Rdb struct {
	Id        uuid.UUID `gorm:"primarykey" json:"id"`
	ProjectId uuid.UUID `json:"project_id"`

	Project    project.Project
	Collection string         `json:"collection"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p Rdb) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var projects []Rdb

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &projects, make([]string, 0))

	y := make([]interface{}, len(projects))
	for i, v := range projects {
		y[i] = v
	}

	return y, nil
}

func (p Rdb) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p Rdb) Total() *int64 {
	return db.MetaDb.TotalRecords(&Rdb{})
}

func (p Rdb) GetById(id string) (interface{}, error) {
	var source Rdb
	conn := db.MetaDb.GetConnection()
	tx := conn.First(&source, "id = ?", id)

	if tx.RowsAffected < 1 {
		return source, errors.New("no found")
	}

	return source, nil
}

func (p Rdb) GetByCollection(collection string) Rdb {
	var source Rdb
	conn := db.MetaDb.GetConnection()
	conn.Preload("Project").First(&source, "collection = ?", collection)
	return source
}

// TableName Gorm table name
func (p Rdb) TableName() string {
	return "rdb"
}
