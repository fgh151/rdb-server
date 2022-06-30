package models

import (
	"db-server/modules/user"
	"db-server/server/db"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type PushMessage struct {
	Id        uuid.UUID      `gorm:"primarykey" json:"id"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Payload   datatypes.JSON `json:"payload"`
	Topic     string         `json:"topic"`
	Sent      bool           `json:"sent"`
	SentAt    time.Time      `json:"sent_at"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Receivers []user.UserDevice `gorm:"many2many:push_receiver;"`
}

type PushLog struct {
	Id            uuid.UUID `gorm:"primarykey"`
	PushMessageId uuid.UUID
	UserDeviceId  uuid.UUID
	Success       bool
	Error         string
	SentAt        time.Time
}

// TableName Gorm table name
func (p PushLog) TableName() string {
	return "push_log"
}

type Sender interface {
	SendPush(message PushMessage, device user.UserDevice) error
}

func (p PushMessage) List(limit int, offset int, sort string, order string, filter map[string]string) []interface{} {
	var pushMessages []PushMessage

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &pushMessages)

	y := make([]interface{}, len(pushMessages))
	for i, v := range pushMessages {
		y[i] = v
	}

	return y
}

func (p PushMessage) GetById(id string) interface{} {
	var pushMessage PushMessage

	conn := db.MetaDb.GetConnection()

	conn.First(&pushMessage, "id = ?", id)

	return pushMessage
}

func (p PushMessage) Delete(id string) {
	if p.Sent == false {
		conn := db.MetaDb.GetConnection()
		conn.Where("id = ?", id).Delete(&p)
	}
}

func (p PushMessage) Total() *int64 {
	return db.MetaDb.TotalRecords(&PushMessage{})
}

// TableName Gorm table name
func (p PushMessage) TableName() string {
	return "push_message"
}
