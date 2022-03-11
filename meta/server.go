package meta

import (
	"db-server/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Connection struct {
	db *gorm.DB
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

func (c Connection) GetConnection() *gorm.DB {

	if c.db == nil {

		c.db, _ = c.connect()

		err := c.db.AutoMigrate(&models.Project{})
		err = c.db.AutoMigrate(&models.User{})
		if err != nil {
			panic("failed to migrate meta database")
		}
	}

	return c.db
}

func (c Connection) GetKey(topic string) string {
	var project models.Project

	conn := c.GetConnection()

	conn.First(&project, "topic = ?", topic)

	return project.Key
}

func (c Connection) GetProjectById(id string) models.Project {
	var project models.Project

	conn := c.GetConnection()

	conn.First(&project, "id = ?", id)

	return project
}

func (c Connection) GetUserById(id string) models.User {
	var user models.User

	conn := c.GetConnection()

	conn.First(&user, "id = ?", id)

	return user
}

func (c Connection) DeleteProjectById(id string) {
	var project models.Project
	conn := c.GetConnection()
	conn.Where("id = ?", id).Delete(&project)
}

func (c Connection) DeleteUserById(id string) {
	var user models.User
	conn := c.GetConnection()
	conn.Where("id = ?", id).Delete(&user)
}

func (c Connection) ListProjects() []models.Project {

	var projects []models.Project

	conn := c.GetConnection()

	conn.Find(&projects)

	return projects
}

func (c Connection) ListUsers() []models.User {

	var users []models.User

	conn := c.GetConnection()

	conn.Find(&users)

	return users
}

var MetaDb = Connection{}
