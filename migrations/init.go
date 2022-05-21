package migrations

import (
	err2 "db-server/err"
	"db-server/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Project{}, &models.User{}, &models.Config{}, &models.DataSource{}, &models.DataSourceEndpoint{}, models.CloudFunction{})

	err2.PanicErr(err)
}
