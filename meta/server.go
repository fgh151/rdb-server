package meta

import (
	"crypto/md5"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"time"
)

type Connection struct {
	db *gorm.DB
}

type Project struct {
	Id        uint           `gorm:"primarykey" json:"id"`
	Topic     string         `json:"topic"`
	Key       string         `json:"key"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Id           uint           `gorm:"primarykey" json:"id"`
	Email        string         `gorm:"index" json:"email"`
	Token        string         `gorm:"index" json:"token"`
	PasswordHash string         `json:"-"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u User) ValidatePassword(password string) bool {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return u.PasswordHash == hash
}

func HashPassword(password string) string {

	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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

		err := c.db.AutoMigrate(&Project{})
		err = c.db.AutoMigrate(&User{})
		if err != nil {
			panic("failed to migrate meta database")
		}
	}

	return c.db
}

func (c Connection) GetKey(topic string) string {
	var project Project

	conn := c.GetConnection()

	conn.First(&project, "topic = ?", topic)

	return project.Key
}

func (c Connection) GetProjectById(id string) Project {
	var project Project

	conn := c.GetConnection()

	conn.First(&project, "id = ?", id)

	return project
}

func (c Connection) GetUserById(id string) User {
	var user User

	conn := c.GetConnection()

	conn.First(&user, "id = ?", id)

	return user
}

func (c Connection) DeleteProjectById(id string) {
	var project Project
	conn := c.GetConnection()
	conn.Where("id = ?", id).Delete(&project)
}

func (c Connection) DeleteUserById(id string) {
	var user User
	conn := c.GetConnection()
	conn.Where("id = ?", id).Delete(&user)
}

func (c Connection) ListProjects() []Project {

	var projects []Project

	conn := c.GetConnection()

	conn.Find(&projects)

	return projects
}

func (c Connection) ListUsers() []User {

	var users []User

	conn := c.GetConnection()

	conn.Find(&users)

	return users
}

var MetaDb = Connection{}
