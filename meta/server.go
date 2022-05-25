package meta

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type Connection struct {
	db *gorm.DB
}

func (c Connection) connect() (*gorm.DB, error) {
	log.Debug("Set meta db driver " + os.Getenv("META_DB_TYPE"))
	switch os.Getenv("META_DB_TYPE") {
	case "sqlite":
		return gorm.Open(sqlite.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	case "mysql":
		return gorm.Open(mysql.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	case "postgres":
		return gorm.Open(postgres.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	}
	panic("failed to connect database")
}

func (c Connection) GetConnection() *gorm.DB {

	if c.db == nil {
		c.db, _ = c.connect()

		switch log.GetLevel() {
		case log.PanicLevel:
			c.db.Logger.LogMode(logger.Silent)
			break
		case log.FatalLevel:
			c.db.Logger.LogMode(logger.Silent)
			break
		case log.ErrorLevel:
			c.db.Logger.LogMode(logger.Error)
			break
		case log.WarnLevel:
			c.db.Logger.LogMode(logger.Warn)
			break
		case log.InfoLevel:
		case log.DebugLevel:
		case log.TraceLevel:
			c.db.Logger.LogMode(logger.Info)
			break
		}
	}

	return c.db
}

var MetaDb = Connection{}
