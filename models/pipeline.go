package models

import (
	"db-server/server"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
	"time"
)

type PipelineOutputType string

const (
	TopicOutput PipelineOutputType = "topic"
)

type PipelineInputType string

const (
	FunctionInput PipelineInputType = "func"
)

type Pipeline struct {
	// The pipeline UUID
	// example: 6204037c-30e6-418b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// Mnemonic name
	Title string `json:"title"`
	// Input type
	Input PipelineInputType `json:"input"`
	// The pipeline input UUID
	// example: 6204037c-30e6-418b-8saa-dd8219860b4b
	InputId uuid.UUID `json:"input_id"`
	// Output type
	Output PipelineOutputType `json:"output"`
	// The pipeline output UUID
	// example: 6204037c-30e6-413b-8saa-dd8219860b4b
	OutputId  uuid.UUID      `json:"output_id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PipelineProcess interface {
	PipelineProcess(data interface{})
}

func (p Pipeline) List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{} {
	var sources []Pipeline

	conn := server.MetaDb.GetConnection()

	log.Debug(filter)

	conn.Limit(limit).Offset(offset).Order(sort + " " + order).Where(filter).Find(&sources)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func (p Pipeline) GetById(id string) interface{} {
	var source Pipeline
	conn := server.MetaDb.GetConnection()
	conn.First(&source, "id = ?", id)
	return source
}

func (p Pipeline) Delete(id string) {
	conn := server.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p Pipeline) Total() *int64 {
	return TotalRecords(&Pipeline{})
}

func RunPipeline(inputName string, inputID uuid.UUID, data interface{}) {
	var source Pipeline
	conn := server.MetaDb.GetConnection()
	tx := conn.First(&source, "input = ? AND input_id = ?", inputName, inputID.String())
	if tx.RowsAffected > 0 {
		switch source.Output {
		case TopicOutput:
			t := Project{}.GetById(source.OutputId.String()).(Project)
			_ = server.SaveTopicMessage(os.Getenv("DB_NAME"), t.Topic, data)
		}
	}
	return
}
