package models

import (
	"db-server/security"
	"db-server/server"
	"github.com/google/uuid"
)

type Model interface {
	List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{}

	GetById(id string) interface{}

	Delete(id string)

	Total() *int64
}

func CreateDemo() {
	var u = User{
		Email:        "test@example.com",
		PasswordHash: security.HashPassword("test"),
		Token:        "123",
		Active:       true,
		Admin:        true,
	}

	u.Id, _ = uuid.NewUUID()

	server.MetaDb.GetConnection().Create(&u)
}

func TotalRecords(m Model) *int64 {
	conn := server.MetaDb.GetConnection()
	var cnt int64
	conn.Model(&m).Count(&cnt)
	return &cnt
}
