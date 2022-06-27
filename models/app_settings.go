package models

import (
	"db-server/server"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AppSettings struct {
	// The UUID
	// example: 6204037c-30e6-408b-8aaa-dd8219860e4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// Mnemonic name
	// example: oauth_gh_client_secret
	Name string `gorm:"index"`
	// Settings value
	// example: 123
	Value string
	// The project UUID
	// example: 6204037c-30e6-438b-8aaa-dd8219860e4b
	ProjectId uuid.UUID `json:"project_id"`
	// The project
	Project   Project        `json:"project"`
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
