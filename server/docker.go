package server

import (
	"context"
	err2 "db-server/err"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

func BuildDockerImage(tar io.Reader, tags []string) error {
	cli, err := GetDockerCli()
	err2.DebugErr(err)

	ctx := context.Background()

	//tar, err := archive.TarWithOptions("node-hello/", &archive.TarOptions{})
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       tags,
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

func PullDockerImage(refStr string) {

	ctx := context.Background()
	cli, err := GetDockerCli()
	if err != nil {
		log.Debug(err)
		return
	}
	reader, err := cli.ImagePull(ctx, refStr, types.ImagePullOptions{})
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	log.Debug(buf.String())
}

func CreateDockerContainer(image string, cmd []string, env []string) (string, error) {
	ctx := context.Background()
	cli, err := GetDockerCli()

	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Env:   env,   //strings.Split(p.Env, "\n"),
		Image: image, //uri.Image,
		Cmd:   cmd,   //prepareDockerParams(p.Params),
	}, nil, nil, nil, "")

	return resp.ID, nil
}

func GetDockerCli() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
