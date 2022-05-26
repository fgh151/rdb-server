package models

import (
	"db-server/messages"
	"db-server/meta"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Receivers []UserDevice `gorm:"many2many:push_receiver;"`
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
		log.Debug("Send push " + p.Id.String() + " to " + receiver.Id.String())
		switch receiver.Device {
		case "ios":
			messages.Ios{}.SendPush()
			break

		case "android":
			messages.Android{}.SendPush()
			break
		}
	}

	p.Sent = true
	meta.MetaDb.GetConnection().Save(&p)
}

type UserDevice struct {
	Id        uuid.UUID      `gorm:"primarykey" json:"id"`
	UserId    uuid.UUID      `json:"user_id"`
	User      User           `json:"-"`
	Device    string         `json:"device"`
	DeviceId  string         `json:"device_id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
