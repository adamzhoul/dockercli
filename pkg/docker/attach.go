package docker

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	kubetype "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/kubelet/dockershim/libdocker"
)

type ContainerAttacher struct {
	client *dockerclient.Client
}

func NewContainerAttacher() *ContainerAttacher {

	return &ContainerAttacher{
		client: client,
	}
}

// attach to container
// 1. handle size
// 2. attach
// 3. hold conn
func (a *ContainerAttacher) AttachContainer(name string, uid kubetype.UID, container string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize) error {

	// handle size
	handleResizing(resize, a.client, container, resizeContainer)

	// attach to container
	opts := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stderr: true,
		Stdout: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	resp, er := a.client.ContainerAttach(ctx, container, opts)
	if er != nil {
		log.Println(er)
		return er
	}
	defer resp.Close()

	// hold attach conn
	sopts := libdocker.StreamOptions{
		InputStream:  in,
		OutputStream: out,
		ErrorStream:  err,
		RawTerminal:  true,
	}
	er = holdHijackedConnection(sopts.RawTerminal, sopts.InputStream, sopts.OutputStream, sopts.ErrorStream, resp)
	if er != nil {
		log.Println(er)
		return er
	}
	return nil
}
