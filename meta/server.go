package meta

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Connection struct {
	db *gorm.DB
}

type Project struct {
	gorm.Model
	Topic string
	Key   string
}

func (c Connection) connect() (*gorm.DB, error) {
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

func (c Connection) getConnection() *gorm.DB {

	if c.db == nil {

		c.db, _ = c.connect()

		err := c.db.AutoMigrate(&Project{})
		if err != nil {
			panic("failed to migrate meta database")
		}
	}

	return c.db
}

func (c Connection) GetKey(topic string) string {
	var project Project

	conn := c.getConnection()

	conn.First(&project, "topic = ?", topic)

	return project.Key
}

var MetaDb = Connection{}
