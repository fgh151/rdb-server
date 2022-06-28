package device

import (
	"db-server/events"
	"db-server/modules/push/models"
	"db-server/modules/user"
)

type InnerPush struct {
}

func (p InnerPush) SendPush(message models.PushMessage, device user.UserDevice) error {
	return events.GetPush().Send(device.DeviceToken, message)
}
