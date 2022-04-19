package migrations

import (
	err2 "db-server/err"
	"db-server/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Project{}, &models.User{})

	err2.CheckErr(err)
}
