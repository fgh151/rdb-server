package device

import (
	err2 "db-server/err"
	"db-server/modules/push/models"
	"db-server/modules/user"
	"encoding/json"
	"fmt"
	"github.com/appleboy/go-fcm"
	log "github.com/sirupsen/logrus"
	"os"
)

type Android struct {
}

func (p Android) SendPush(message models.PushMessage, device user.UserDevice) error {

	log.Debug("Send push " + message.Id.String() + " to " + device.Id.String())

	var data map[string]interface{}

	err := json.Unmarshal([]byte(message.Body), &data)
	err2.DebugErr(err)
	if err != nil {
		return err
	}

	msg := &fcm.Message{
		To:   device.DeviceToken,
		Data: data,
		Notification: &fcm.Notification{
			Title: message.Title,
			Body:  message.Body,
		},
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient(os.Getenv("PUSH_FCM_API_KEY"))
	err2.DebugErr(err)
	if err != nil {
		return err
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	err2.DebugErr(err)
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("%#v\n", response))
	return nil
}
