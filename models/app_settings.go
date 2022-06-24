package models

import (
	"db-server/server"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AppSettings struct {
	Id        uint   `gorm:"primaryKey"`
	Name      string `gorm:"index"`
	Value     string
	ProjectId uuid.UUID `json:"project_id"`
	Project   Project   `json:"project"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func GetAppSettingsByName(projectId uuid.UUID, name string) string {
	s := AppSettings{}
	conn := server.MetaDb.GetConnection()
	tx := conn.First(&s, "name = ? AND project_id = ? ", name, projectId)
	if tx.RowsAffected < 1 {
		return ""
	}
	return s.Value
}
