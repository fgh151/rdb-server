package migrations

import (
	err2 "db-server/err"
	"db-server/messages"
	"db-server/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Project{},
		&models.User{},
		&models.Config{},
		&models.DataSource{},
		&models.DataSourceEndpoint{},
		&models.CloudFunction{},
		&models.CloudFunctionLog{},
		&messages.PushMessage{},
		&messages.UserDevice{},
	)

	err2.PanicErr(err)
}
