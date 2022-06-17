package server

import (
	"db-server/drivers"
	"db-server/events"
)

// SaveTopicMessage
// Save document to db and register new message
func SaveTopicMessage(db string, topic string, payload interface{}) error {
	_, err := drivers.GetDbInstance().Insert(db, topic, payload)
	if err == nil {
		events.GetInstance().RegisterNewMessage(topic, payload)
	}

	return err
}
