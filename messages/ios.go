package messages

import (
	err2 "db-server/err"
	"fmt"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	log "github.com/sirupsen/logrus"
	"os"
)

type Ios struct {
}

func (p Ios) SendPush(message PushMessage, device UserDevice) {

	log.Debug("Send push " + message.Id.String() + " to " + device.Id.String())

	cert, pemErr := certificate.FromPemFile(os.Getenv("PUSH_APNS_PEM_FILE"), os.Getenv("PUSH_APNS_PEM_FILE_PASSWORD"))
	err2.DebugErr(pemErr)

	//payload := NewPayload().Alert("hello").Badge(1).Custom("key", "val")

	notification := &apns.Notification{}
	notification.DeviceToken = device.DeviceToken
	notification.Topic = message.Topic
	notification.Payload = []byte(message.Payload) // See Payload section below

	client := apns.NewClient(cert).Development()
	response, err := client.Push(notification)
	err2.DebugErr(err)

	fmt.Println(response)

	log.Debug(fmt.Sprintf("%#v\n", response))
}
