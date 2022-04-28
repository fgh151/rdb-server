package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type DataSource struct {
	Id  uuid.UUID `gorm:"primarykey" json:"id"`
	Dsn string    `json:"dsn"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
