package plugin

import (
	"db-server/server/db"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Plugin struct {
	// The plugin UUID
	// example: 6204011c-30s6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`

	// Plugin filename
	FileName  string         `json:"file_name"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName Gorm table name
func (p Plugin) TableName() string {
	return "plugin"
}

func (p Plugin) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var models []Plugin

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &models, make([]string, 0))

	y := make([]interface{}, len(models))
	for i, v := range models {
		y[i] = v
	}

	return y, nil
}

func (p Plugin) Total() *int64 {
	return db.MetaDb.TotalRecords(&Plugin{})
}

func (p Plugin) GetById(id string) (interface{}, error) {
	var source Plugin

	conn := db.MetaDb.GetConnection()

	tx := conn.First(&source, "id = ?", id)

	if tx.RowsAffected < 1 {
		return source, errors.New("no found")
	}

	return source, nil
}

func (p Plugin) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p Plugin) Run(params interface{}) {

}
