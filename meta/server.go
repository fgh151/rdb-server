package meta

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Connection struct {
	db *gorm.DB
}

func (c Connection) connect() *gorm.DB {
	log.Debug("Set meta db driver " + os.Getenv("META_DB_TYPE"))

	var db *gorm.DB
	var err error

	switch os.Getenv("META_DB_TYPE") {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(os.Getenv("META_DB_DSN")), &gorm.Config{})
	}

	if err != nil {
		panic("failed to connect database")
	}

	if log.GetLevel() >= log.DebugLevel {
		return db.Debug()
	}

	return db
}

func (c Connection) GetConnection() *gorm.DB {

	if c.db == nil {
		c.db = c.connect()
	}

	return c.db
}

var MetaDb = Connection{}
