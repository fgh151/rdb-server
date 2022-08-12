package migrations

import (
	err2 "db-server/err"
	"db-server/modules/cf"
	"db-server/modules/config"
	"db-server/modules/cron"
	"db-server/modules/ds"
	"db-server/modules/oauth"
	"db-server/modules/pipeline"
	"db-server/modules/plugin"
	"db-server/modules/project"
	"db-server/modules/push/models"
	"db-server/modules/rdb"
	"db-server/modules/settings"
	"db-server/modules/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&project.Project{},
		&user.User{},
		&config.Config{},
		&ds.DataSource{},
		&ds.DataSourceEndpoint{},
		&cf.CloudFunction{},
		&cf.CloudFunctionLog{},
		&models.PushMessage{},
		&user.UserDevice{},
		&models.PushLog{},
		&cron.CronJob{},
		&pipeline.Pipeline{},
		&oauth.UserOauth{},
		&settings.AppSettings{},
		&rdb.Rdb{},
		&plugin.Plugin{},
	)

	err2.PanicErr(err)
}
