package models

import "db-server/events"

type InnerPush struct {
}

func (p InnerPush) SendPush(message PushMessage, device UserDevice) error {
	return events.GetPush().Send(device.DeviceToken, message)
}
