package events

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type subscribers struct {
	list map[string][]*websocket.Conn
}

var Subscribers = subscribers{}

func init() {
	Subscribers.list = make(map[string][]*websocket.Conn)
}

func Subscribe(topic string, listener *websocket.Conn) {

	var currentList []*websocket.Conn

	if _, ok := Subscribers.list[topic]; ok {
	} else {
		currentList = []*websocket.Conn{}
	}

	currentList = append(currentList, listener)

	Subscribers.list[topic] = currentList
}

func Unsubscribe(topic string, listener *websocket.Conn) {
	if _, ok := Subscribers.list[topic]; ok {
		for i, val := range Subscribers.list[topic] {
			if val == listener {
				remove(Subscribers.list[topic], i)
			}
		}
	}
}

func remove(s []*websocket.Conn, i int) []*websocket.Conn {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RegisterNewMessage(topic string, content interface{}) {
	if currentList, ok := Subscribers.list[topic]; ok {
		for _, listener := range currentList {
			msg, _ := json.Marshal(content)
			listener.WriteMessage(1, msg)
		}
	}
}
