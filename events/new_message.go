package events

import (
	err2 "db-server/err"
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type subscribers struct {
	list map[string][]*websocket.Conn
}

type EventHandler struct {
	subscribers subscribers
	sync.RWMutex
}

func (e *EventHandler) Subscribe(topic string, listener *websocket.Conn) {

	var currentList []*websocket.Conn

	if _, ok := e.subscribers.list[topic]; ok {
	} else {
		currentList = []*websocket.Conn{}
	}

	e.Lock()
	defer e.Unlock()

	currentList = append(currentList, listener)

	e.subscribers.list[topic] = currentList
}

func (e *EventHandler) Unsubscribe(topic string, listener *websocket.Conn) {
	if _, ok := e.subscribers.list[topic]; ok {
		for i, val := range e.subscribers.list[topic] {
			if val == listener {
				e.Lock()
				e.remove(e.subscribers.list[topic], i)
				e.Unlock()
			}
		}
	}
}

func (e *EventHandler) remove(s []*websocket.Conn, i int) []*websocket.Conn {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (e *EventHandler) RegisterNewMessage(topic string, content interface{}) {
	if currentList, ok := e.subscribers.list[topic]; ok {
		for _, listener := range currentList {
			msg, _ := json.Marshal(content)
			err := listener.WriteMessage(1, msg)
			err2.DebugErr(err)
		}
	}
}

var instance *EventHandler = nil

func GetInstance() *EventHandler {
	if instance == nil {
		instance = new(EventHandler)
		instance.subscribers.list = make(map[string][]*websocket.Conn)
	}
	return instance
}
