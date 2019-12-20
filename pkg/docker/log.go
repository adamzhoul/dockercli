package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func TailLog(container string) error {

	if !strings.HasPrefix(container, dockerContainerPrefix) {
		return errors.New(fmt.Sprintf("not docker container:%s", container))
	}

	dockerContainerId := container[len(dockerContainerPrefix):]
	log.Println("exec attach:", dockerContainerId)

	resp, err := client.ContainerLogs(context.Background(), dockerContainerId, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true})
	if err != nil {
		return err
	}

	//resp.Read()
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return err
	}

	return nil
}
