package docker

import (
	"context"
	"io"
	"os"

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

	// handle size
	handleResizing(resize, e.client, container, resizeContainer)

	// execCreate
	respIdExecCreate, er := e.client.ContainerExecCreate(context.Background(), container,
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

	// attach
	resp, er := e.client.ContainerExecAttach(context.Background(), respIdExecCreate.ID, types.ExecStartCheck{
		Tty: true,
	})
	if er != nil {
		return er
	}
	defer resp.Close()

	// hold hijack
	er = holdHijackedConnection(true, os.Stdin, os.Stdout, os.Stderr, resp)
	if er != nil {
		return er
	}

	return nil
}
