package models

import (
	"db-server/server"
	"gorm.io/gorm"
	"time"
)

type AppSettings struct {
	Id    uint   `gorm:"primaryKey"`
	Name  string `gorm:"index"`
	Value string

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func GetAppSettingsByName(name string) string {
	s := AppSettings{}
	conn := server.MetaDb.GetConnection()
	tx := conn.First(&s, "name = ?", name)
	if tx.RowsAffected < 1 {
		return ""
	}
	return s.Value
}
