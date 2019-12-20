package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	kubetype "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"
)

type ContainerLogAttacher struct {
	client *dockerclient.Client
}

func NewContainerLogAttacher() *ContainerLogAttacher {

	return &ContainerLogAttacher{
		client: client,
	}
}
func (l *ContainerLogAttacher) AttachContainer(name string, uid kubetype.UID, container string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize) error {

	if !strings.HasPrefix(container, dockerContainerPrefix) {
		return errors.New(fmt.Sprintf("not docker container:%s", container))
	}

	dockerContainerId := container[len(dockerContainerPrefix):]
	log.Println("exec attach:", dockerContainerId)

	resp, er := l.client.ContainerLogs(context.Background(), dockerContainerId, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true})
	if er != nil {
		return er
	}

	_, er = io.Copy(out, resp)
	if er != nil {
		return er
	}

	return nil
}
