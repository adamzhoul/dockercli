package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	//"os"
	"strings"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	kubetype "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"
)

type ContainerExecAttacher struct {
	client *dockerclient.Client
}

func NewContainerExecAttacher() *ContainerExecAttacher {

	return &ContainerExecAttacher{
		client: client,
	}
}

func (e *ContainerExecAttacher) AttachContainer(name string, uid kubetype.UID, container string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize) error {

	if !strings.HasPrefix(container, dockerContainerPrefix) {
		return errors.New(fmt.Sprintf("not docker container:%s", container))
	}

	dockerContainerId := container[len(dockerContainerPrefix):]
	log.Println("exec attach:", dockerContainerId)

	// execCreate
	respIdExecCreate, er := e.client.ContainerExecCreate(context.Background(), dockerContainerId,
		types.ExecConfig{
			//	User:         "1000",
			Tty:          true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          []string{"/bin/bash"},
		})
	if er != nil {
		return er
	}

	// handle size
	handleResizing(resize, e.client, respIdExecCreate.ID, resizeExecContainer)

	// attach
	ctx, cancel := context.WithCancel(context.Background())
	resp, er := e.client.ContainerExecAttach(ctx, respIdExecCreate.ID, types.ExecStartCheck{
		Tty: true,
	})
	defer cancel() // 关闭/bin/bash进程
	if er != nil {
		return er
	}
	defer resp.Close()

	// hold hijack
	er = holdHijackedConnection(true, in, out, out, resp)
	if er != nil {
		return er
	}

	return nil
}
