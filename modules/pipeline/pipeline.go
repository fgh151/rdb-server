package pipeline

import (
	"db-server/modules/plugin"
	"db-server/modules/rdb"
	"db-server/server"
	"db-server/server/db"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"time"
)

type PipelineOutputType string

const (
	TopicOutput  PipelineOutputType = "topic"
	PluginOutput PipelineOutputType = "plugin"
)

type PipelineInputType string

//const (
//	FunctionInput PipelineInputType = "func"
//)

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

// TableName Gorm table name
func (p Pipeline) TableName() string {
	return "pipeline"
}

type PipelineProcess interface {
	PipelineProcess(data interface{})
}

func (p Pipeline) List(limit int, offset int, sort string, order string, filter map[string]string) ([]interface{}, error) {
	var sources []Pipeline

	db.MetaDb.ListQuery(limit, offset, sort, order, filter, &sources, make([]string, 0))

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y, nil
}

func (p Pipeline) GetById(id string) (interface{}, error) {
	var source Pipeline
	conn := db.MetaDb.GetConnection()
	tx := conn.First(&source, "id = ?", id)

	if tx.RowsAffected < 1 {
		return source, errors.New("no found")
	}

	return source, nil
}

func (p Pipeline) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p Pipeline) Total() *int64 {
	return db.MetaDb.TotalRecords(&Pipeline{})
}

func RunPipeline(inputName string, inputID uuid.UUID, data interface{}) {
	var source Pipeline
	conn := db.MetaDb.GetConnection()
	tx := conn.First(&source, "input = ? AND input_id = ?", inputName, inputID.String())
	if tx.RowsAffected > 0 {
		switch source.Output {
		case TopicOutput:
			t, err := rdb.Rdb{}.GetById(source.OutputId.String())
			if err == nil {
				_ = server.SaveTopicMessage(os.Getenv("DB_NAME"), t.(rdb.Rdb).Collection, data)
			}
		case PluginOutput:
			p, err := plugin.Plugin{}.GetById(source.OutputId.String())
			if err == nil {
				p.(plugin.Plugin).Run(data)
			}
		}
	}
	return
}
