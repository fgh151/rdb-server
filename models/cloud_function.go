package models

import (
	"context"
	err2 "db-server/err"
	"db-server/server"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"strings"
	"time"
)

//  cf := models.CloudFunction{
//		Container: "docker.io/library/alpine",
//		Params:    []string{"echo", "hello world"},
//	}
//
//	cf.Run()

type CloudFunction struct {
	Id        uuid.UUID      `gorm:"primarykey" json:"id"`
	ProjectId uuid.UUID      `json:"project_id"`
	Title     string         `json:"title"`
	Container string         `json:"container"`
	Params    string         `json:"params"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Project   Project

	RunCount int64 `gorm:"-:all" json:"run_count"`
}

type CloudFunctionLog struct {
	Id         uuid.UUID `gorm:"primarykey" json:"id"`
	FunctionId uuid.UUID `json:"function_id"`
	RunAt      time.Time `json:"run_at"`
	Result     string    `json:"result"`
}

func ListCfLog(fId uuid.UUID, limit int, offset int, sort string, order string) []interface{} {
	var sources []CloudFunctionLog

	conn := server.MetaDb.GetConnection()

	conn.Find(&sources, CloudFunctionLog{FunctionId: fId}).Limit(limit).Offset(offset).Order(order + " " + sort)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func LogsTotal(fId uuid.UUID) *int64 {
	conn := server.MetaDb.GetConnection()
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

func (p CloudFunction) List(limit int, offset int, sort string, order string) []interface{} {
	var sources []CloudFunction

	conn := server.MetaDb.GetConnection()

	conn.Limit(limit).Offset(offset).Order(order + " " + sort).Find(&sources)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		v.RunCount = *LogsTotal(v.Id)
		y[i] = v
	}

	return y
}

func (p CloudFunction) Total() *int64 {
	return TotalRecords(&CloudFunction{})
}

func (p CloudFunction) GetById(id string) interface{} {
	var source CloudFunction

	conn := server.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	uid, _ := uuid.Parse(id)

	source.RunCount = *LogsTotal(uid)

	return source
}

func (p CloudFunction) Delete(id string) {
	conn := server.MetaDb.GetConnection()
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

	uri, err := GetContainerUri(p.Container)

	err2.WarnErr(err)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	p.checkErr(runId, err)

	path := uri.Host + "/" + uri.Vendor + "/" + uri.Image

	// Делаем docker pull
	reader, err := cli.ImagePull(ctx, path, types.ImagePullOptions{})
	p.checkErr(runId, err)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	log.Debug(buf.String())

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: uri.Image,
		Cmd:   prepareDockerParams(p.Params),
	}, nil, nil, nil, "")
	p.checkErr(runId, err)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		err2.DebugErr(err)
		p.log(runId, "error "+err.Error())
		return
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		p.checkErr(runId, err)
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: log.GetLevel() >= log.InfoLevel})
	p.checkErr(runId, err)

	result, err := makeResultFromStream(out)
	p.checkErr(runId, err)

	log.Debug("Cf run result " + runId.String() + " " + result)

	p.log(runId, result)
}

func BuildImage(tar io.Reader, uri ContainerUri) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err2.DebugErr(err)

	ctx := context.Background()

	//tar, err := archive.TarWithOptions("node-hello/", &archive.TarOptions{})
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{uri.Vendor + "/" + uri.Image},
		Remove:     true,
	}

	res, err := cli.ImageBuild(ctx, tar, opts)

	defer func() {
		err := res.Body.Close()
		err2.DebugErr(err)
	}()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)

	if err != nil {
		log.Debug(err)
		return err
	}

	// check errors
	fmt.Println(buf.String())

	if buf.String() != "" {
		return errors.New(buf.String())
	}

	return nil
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
	server.MetaDb.GetConnection().Create(&flog)
}
