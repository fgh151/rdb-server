package events

import (
	err2 "db-server/err"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type pushSubscribers struct {
	list map[string]*websocket.Conn
}

type PushHandler struct {
	subscribers pushSubscribers
	sync.RWMutex
}

func (e *PushHandler) Subscribe(deviceId string, listener *websocket.Conn) {
	e.subscribers.list[deviceId] = listener
}

func (e *PushHandler) Unsubscribe(deviceId string) {
	delete(e.subscribers.list, deviceId)
}

func (e *PushHandler) remove(s []*websocket.Conn, i int) []*websocket.Conn {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (e *PushHandler) RegisterNewMessage(content interface{}) {
	for _, listener := range e.subscribers.list {
		msg, _ := json.Marshal(content)
		err := listener.WriteMessage(websocket.TextMessage, msg)
		err2.DebugErr(err)
	}
}

func (e *PushHandler) Send(deviceId string, payload interface{}) error {

	if conn, ok := e.subscribers.list[deviceId]; ok {

		b, err := json.Marshal(payload)

		err2.DebugErr(err)

		if err == nil {
			err = conn.WriteMessage(websocket.TextMessage, b)
			err2.DebugErr(err)
		}

		return err
	}
	return errors.New("Disconected device")
}

var pushInstance *PushHandler = nil

func GetPush() *PushHandler {
	if pushInstance == nil {
		pushInstance = new(PushHandler)
		pushInstance.subscribers.list = make(map[string]*websocket.Conn)
	}
	return pushInstance
}
