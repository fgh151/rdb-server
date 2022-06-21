package migrations

import (
	err2 "db-server/err"
	"db-server/models"
	"db-server/oauth"
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
		&models.PushMessage{},
		&models.UserDevice{},
		&models.PushLog{},
		&models.CronJob{},
		&models.Pipeline{},
		&oauth.UserOauth{},
		&models.AppSettings{},
	)

	err2.PanicErr(err)
}
