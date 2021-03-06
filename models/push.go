package models

import (
	"db-server/server"
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

	Receivers []UserDevice `gorm:"many2many:push_receiver;"`
}

type PushLog struct {
	Id            uuid.UUID `gorm:"primarykey"`
	PushMessageId uuid.UUID
	UserDeviceId  uuid.UUID
	Success       bool
	Error         string
	SentAt        time.Time
}

type Sender interface {
	SendPush(message PushMessage, device UserDevice) error
}

func (p PushMessage) List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{} {
	var pushMessages []PushMessage

	conn := server.MetaDb.GetConnection()

	conn.Offset(offset).Limit(limit).Order(sort + " " + order).Where(filter).Find(&pushMessages)

	y := make([]interface{}, len(pushMessages))
	for i, v := range pushMessages {
		y[i] = v
	}

	return y
}

func (p PushMessage) GetById(id string) interface{} {
	var pushMessage PushMessage

	conn := server.MetaDb.GetConnection()

	conn.First(&pushMessage, "id = ?", id)

	return pushMessage
}

func (p PushMessage) Delete(id string) {
	if p.Sent == false {
		conn := server.MetaDb.GetConnection()
		conn.Where("id = ?", id).Delete(&p)
	}
}

func (p PushMessage) Total() *int64 {
	return TotalRecords(&PushMessage{})
}

func (p PushMessage) Send() {

	for _, receiver := range p.Receivers {
		switch receiver.Device {
		case "ios":
			createPushLog(
				p,
				receiver,
				Ios{}.SendPush(p, receiver),
			)
			break

		case "android":
			createPushLog(
				p,
				receiver,
				Android{}.SendPush(p, receiver),
			)
			break
		default:
			createPushLog(p, receiver, InnerPush{}.SendPush(p, receiver))
		}
	}

	p.Sent = true
	p.SentAt = time.Now()
	server.MetaDb.GetConnection().Save(&p)
}

func createPushLog(message PushMessage, device UserDevice, err error) {

	id, _ := uuid.NewUUID()
	log := PushLog{
		Id:            id,
		PushMessageId: message.Id,
		UserDeviceId:  device.Id,
		Success:       err == nil,
		Error:         err.Error(),
		SentAt:        time.Now(),
	}

	server.MetaDb.GetConnection().Create(&log)
}

// swagger:model
type UserDevice struct {
	// The UUID of a device
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// The UUID of owned user
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	UserId uuid.UUID `json:"user_id"`
	User   User      `json:"-"`
	// Device type can be android | ios | macos | windows | linux | web
	// example: 'linux'
	Device string `json:"device"`
	// UUID Device Token
	// 9f80fcc8-0102-4795-94c0-4190c168ffc2
	DeviceToken string         `json:"device_token"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
