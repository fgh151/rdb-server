package messages

import (
	"db-server/meta"
	"db-server/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PushMessage struct {
	Id        uuid.UUID      `gorm:"primarykey" json:"id"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Payload   string         `json:"payload"`
	Topic     string         `json:"topic"`
	Sent      bool           `json:"sent"`
	SentAt    time.Time      `json:"sent_at"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Receivers []UserDevice `gorm:"many2many:push_receiver;"`
}

type Sender interface {
	SendPush(message PushMessage, device UserDevice)
}

func (p PushMessage) List(limit int, offset int, sort string, order string) []interface{} {
	var pushMessages []PushMessage

	conn := meta.MetaDb.GetConnection()

	conn.Offset(offset).Limit(limit).Order(order + " " + sort).Find(&pushMessages)

	y := make([]interface{}, len(pushMessages))
	for i, v := range pushMessages {
		y[i] = v
	}

	return y
}

func (p PushMessage) GetById(id string) interface{} {
	var pushMessage PushMessage

	conn := meta.MetaDb.GetConnection()

	conn.First(&pushMessage, "id = ?", id)

	return pushMessage
}

func (p PushMessage) Delete(id string) {
	if p.Sent == false {
		conn := meta.MetaDb.GetConnection()
		conn.Where("id = ?", id).Delete(&p)
	}
}

func (p PushMessage) Total() *int64 {
	conn := meta.MetaDb.GetConnection()
	var pushMessages []PushMessage
	var cnt int64
	conn.Find(&pushMessages).Count(&cnt)

	return &cnt
}

func (p PushMessage) Send() {

	for _, receiver := range p.Receivers {
		switch receiver.Device {
		case "ios":
			Ios{}.SendPush(p, receiver)
			break

		case "android":
			Android{}.SendPush(p, receiver)
			break
		}
	}

	p.Sent = true
	p.SentAt = time.Now()
	meta.MetaDb.GetConnection().Save(&p)
}

type UserDevice struct {
	Id          uuid.UUID      `gorm:"primarykey" json:"id"`
	UserId      uuid.UUID      `json:"user_id"`
	User        models.User    `json:"-"`
	Device      string         `json:"device"`
	DeviceToken string         `json:"device_token"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
