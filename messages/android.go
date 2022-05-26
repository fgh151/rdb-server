package messages

import (
	err2 "db-server/err"
	"fmt"
	"github.com/appleboy/go-fcm"
)

type Android struct {
}

func (p Android) SendPush() {
	msg := &fcm.Message{
		To: "sample_device_token",
		Data: map[string]interface{}{
			"foo": "bar",
		},
		Notification: &fcm.Notification{
			Title: "title",
			Body:  "body",
		},
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient("sample_api_key")
	err2.DebugErr(err)

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	err2.DebugErr(err)

	fmt.Printf("%#v\n", response)
}
