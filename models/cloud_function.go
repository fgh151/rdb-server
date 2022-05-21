package models

import (
	"context"
	"db-server/meta"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"os"
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
}

type containerUri struct {
	Host    string
	Vendor  string
	Image   string
	Version string
}

func getContainerUri(source string) containerUri {

	uri := containerUri{}
	parts := strings.Split(source, ":")

	if len(parts) == 2 {
		uri.Version = parts[1]
	}

	pathParts := strings.Split(parts[0], "/")

	uri.Host = pathParts[0]
	uri.Vendor = pathParts[1]
	uri.Image = pathParts[2]

	return uri
}

func (p CloudFunction) List(limit int, offset int, sort string, order string) []interface{} {
	var sources []CloudFunction

	conn := meta.MetaDb.GetConnection()

	conn.Find(&sources).Limit(limit).Offset(offset).Order(order + " " + sort)

	y := make([]interface{}, len(sources))
	for i, v := range sources {
		y[i] = v
	}

	return y
}

func (p CloudFunction) Total() *int64 {
	conn := meta.MetaDb.GetConnection()
	var sources []CloudFunction
	var cnt int64
	conn.Find(&sources).Count(&cnt)

	return &cnt
}

func (p CloudFunction) GetById(id string) interface{} {
	var source CloudFunction

	conn := meta.MetaDb.GetConnection()

	conn.First(&source, "id = ?", id)

	return source
}

func (p CloudFunction) Delete(id string) {
	conn := meta.MetaDb.GetConnection()
	conn.Where("id = ?", id).Delete(&p)
}

func (p CloudFunction) Run() {

	uri := getContainerUri(p.Container)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	path := uri.Host + "/" + uri.Vendor + "/" + uri.Image

	// Делаем docker pull
	reader, err := cli.ImagePull(ctx, path, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: uri.Image,
		Cmd:   strings.Split(p.Params, "\\"),
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}
