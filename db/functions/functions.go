package functions

import (
	"db-server/db"
	"db-server/models"
	"encoding/json"
)

func Update(id interface{}, value interface{}) error {

	m := models.Message{}.GetById(id.(string)).(models.Message)
	str, err := json.Marshal(value)
	m.Content = string(str)

	db.DB.GetConnection().Save(&m)

	return err
}

func Insert(topic string, value interface{}) (models.Message, error) {

	t := models.Topic{}.GetByTitle(topic).(models.Topic)

	str, err := json.Marshal(value)

	m := models.Message{Content: string(str), TopicId: t.Id}

	db.DB.GetConnection().Create(&m)

	return m, err
}

func Find(topic string, filter map[string]interface{}) []models.Message {
	t := models.Topic{}.GetByTitle(topic).(models.Topic)
	filter["topic_id"] = t.Id
	return models.Message{}.Find(filter)
}

func List(topic string) []models.Message {
	t := models.Topic{}.GetByTitle(topic).(models.Topic)
	return t.Messages
}
