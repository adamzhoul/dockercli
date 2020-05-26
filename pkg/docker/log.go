package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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
	defer resp.Close()

	
	er = holdLogConnection(in, out, out, resp)
	if er != nil {
		return er
	}
	return nil
}

func holdLogConnection(inputStream io.Reader, outputStream, errorStream io.Writer, resp io.ReadCloser) error {

	receiveStdout := make(chan error)
	if outputStream != nil || errorStream != nil {
		go func() {
			_, er := stdcopy.StdCopy(outputStream, errorStream, resp)
			receiveStdout <- er
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		io.Copy(ioutil.Discard, inputStream)
		stdinDone <- struct{}{}
	}()

	select {
	case err := <-receiveStdout:
		return err
	case <-stdinDone:
		log.Println("stdin done")
	}
	return nil
}
