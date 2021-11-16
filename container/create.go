package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/golang-test-task/logger"
	"github.com/sirupsen/logrus"
)

// CreateContainer creates and starts a container with the provided image and cmd and send its logs to cloudwatch
func CreateContainer(image, bashCmd, group, stream, accessKey, secretKey, region string) error {
	err := logger.Init(accessKey, secretKey, region)
	if err != nil {
		logrus.Errorln("unable to initalize aws logging session")
	}
	err = logger.EnsureLogGroupExists(group)
	if err != nil {
		logrus.Errorln("unable to create log group", err)
		return err
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorln("unable to create docker client", err)
		return err
	}
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		logrus.Errorln("unable to pull docker image", err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", bashCmd},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		logrus.Errorln("unable to create container", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		logrus.Errorln("unable to start the container", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			logrus.Errorln("unable to wait for contianer state", err)
		}
	case <-statusCh:
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		logrus.Errorln("unable to fetch the container logs", err)
	}
	defer out.Close()
	p := make([]byte, 1024)

	for {
		n, err := out.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				logger.SendLogsToCloudwatch(group, stream, string(p[:n])) // should handle any remainding bytes.
				break
			}
			logrus.Errorln("unable to read logs", err)

		}
		logger.SendLogsToCloudwatch(group, stream, string(p[:n]))
	}
	defer out.Close()
	return nil
}
