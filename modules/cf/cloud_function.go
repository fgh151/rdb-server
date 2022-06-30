package cf

import (
	"context"
	err2 "db-server/err"
	"db-server/modules/pipeline"
	"db-server/modules/project"
	"db-server/server"
	"db-server/server/db"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"strings"
	"time"
)

// swagger:model
type CloudFunction struct {
	// The function UUID
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// The project UUID
	// example: 6204037c-30e6-403b-8aaa-dd8219860b4b
	ProjectId uuid.UUID `json:"project_id"`
	// Function title
	Title string `json:"title"`
	// Container name
	// example: docker.io/library/alpine
	Container string `json:"container"`
	// Container run params
	// example: echo test
	Params string `json:"params"`
	// Container env variables
	Env         string         `json:"env"`
	ContainerId string         `json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	// Linked project
	Project project.Project
	// Function run count
	RunCount int64 `gorm:"-:all" json:"run_count"`
}

// TableName Gorm table name
func (p CloudFunction) TableName() string {
	return "cf_function"
}

// swagger:model
type CloudFunctionLog struct {
	// The log UUID
	// example: 6204037c-30e6-408b-8aaa-dd8219860b4b
	Id uuid.UUID `gorm:"primarykey" json:"id"`
	// The function UUID
	// example: 6204037c-30e6-408b-8aaa-dd8219520b4b
	FunctionId uuid.UUID `json:"function_id"`
	// Run date time
	RunAt time.Time `json:"run_at"`
	// Run result
	Result string `json:"result"`
}

// TableName Gorm table name
func (p CloudFunctionLog) TableName() string {
	return "cf_log"
}

func ListCfLog(fId uuid.UUID, limit int, offset int, sort string, order string) []interface{} {
	var sources []CloudFunctionLog

	conn := db.MetaDb.GetConnection()

	conn.Find(&sources, CloudFunctionLog{FunctionId: fId}).Limit(limit).Offset(offset).Order(clause.OrderByColumn{Column: clause.Column{Name: sort}, Desc: order != "ASC"})

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func LogsTotal(fId uuid.UUID) *int64 {
	conn := db.MetaDb.GetConnection()
	var sources []CloudFunctionLog
	var cnt int64
	conn.Count(&cnt).Find(&sources, CloudFunctionLog{FunctionId: fId})

	return &cnt
}

type ContainerUri struct {
	Host    string
	Vendor  string
	Image   string
	Version string
}

func GetContainerUri(source string) (ContainerUri, error) {

	log.Debug("Parse uri " + source)
	uri := ContainerUri{}
	parts := strings.Split(source, ":")

	if len(parts) == 2 {
		uri.Version = parts[1]
	}

	pathParts := strings.Split(parts[0], "/")

	if len(pathParts) < 3 {
		return ContainerUri{}, errors.New("Wrong source " + source)
	}

	uri.Host = pathParts[0]
	uri.Vendor = pathParts[1]
	uri.Image = pathParts[2]

	return uri, nil
}

func (p CloudFunction) List(limit int, offset int, sort string, order string, filter map[string]interface{}) []interface{} {
	var sources []CloudFunction

	conn := db.MetaDb.GetConnection()

	log.Debug(filter)

	conn.Limit(limit).Offset(offset).Order(sort + " " + order).Where(filter).Find(&sources)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		v.RunCount = *LogsTotal(v.Id)
		y[i] = v
	}

	return y
}

func (p CloudFunction) Total() *int64 {
	return db.MetaDb.TotalRecords(&CloudFunction{})
}

func (p CloudFunction) GetById(id string) interface{} {
	var source CloudFunction

	conn := db.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	uid, _ := uuid.Parse(id)

	source.RunCount = *LogsTotal(uid)

	return source
}

func (p CloudFunction) Delete(id string) {
	conn := db.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func prepareDockerParams(raw string) []string {
	parts := strings.Split(raw, "\\")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (p CloudFunction) Run(runId uuid.UUID) {

	if p.ContainerId == "" {
		uri, err := GetContainerUri(p.Container)
		err2.WarnErr(err)

		if err != nil {
			return
		}

		cid, err := server.CreateDockerContainer(uri.Image, prepareDockerParams(p.Params), strings.Split(p.Env, "\n"))
		err2.WarnErr(err)

		if err != nil {
			return
		}

		p.ContainerId = cid

		db.MetaDb.GetConnection().Model(CloudFunction{}).Where("id = ?", p.Id).Update("container_id", cid)
	}

	ctx := context.Background()
	cli, err := server.GetDockerCli()

	if err := cli.ContainerStart(ctx, p.ContainerId, types.ContainerStartOptions{}); err != nil {
		err2.DebugErr(err)
		p.log(runId, "error "+err.Error())
		return
	}

	statusCh, errCh := cli.ContainerWait(ctx, p.ContainerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		p.checkErr(runId, err)
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, p.ContainerId, types.ContainerLogsOptions{ShowStdout: log.GetLevel() >= log.InfoLevel})
	p.checkErr(runId, err)

	result, err := makeResultFromStream(out)
	p.checkErr(runId, err)

	pipeline.RunPipeline("func", p.Id, result)

	log.Debug("Cf run result " + runId.String() + " " + result)

	p.log(runId, result)
}

func makeResultFromStream(stream io.Reader) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, stream)

	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(buf.String())

	//see https://pkg.go.dev/github.com/pborman/ansi#pkg-constants
	result = strings.ReplaceAll(result, "\001", "")
	result = strings.ReplaceAll(result, "\000", "")
	result = strings.ReplaceAll(result, "\005", "")

	return result, nil
}

func (p CloudFunction) checkErr(id uuid.UUID, err error) {
	if err != nil {
		err2.DebugErr(err)
		p.log(id, "error "+err.Error())
		return
	}
}

func (p CloudFunction) log(id uuid.UUID, result string) {
	flog := CloudFunctionLog{
		FunctionId: p.Id,
		RunAt:      time.Now(),
		Result:     result,
		Id:         id,
	}
	var err error
	err2.DebugErr(err)
	db.MetaDb.GetConnection().Create(&flog)
}
