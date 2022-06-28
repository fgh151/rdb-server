package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

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

// TableName Gorm table name
func (d UserDevice) TableName() string {
	return "user_device"
}
