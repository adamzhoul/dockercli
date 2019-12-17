package docker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
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

func (a *ContainerAttacher) AttachContainer(name string, uid kubetype.UID, container string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize) error {

	opts := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stderr: true,
		Stdout: true,
	}
	ctx := context.Background()
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
	er = a.holdHijackedConnection(sopts.RawTerminal, sopts.InputStream, sopts.OutputStream, sopts.ErrorStream, resp)
	if er != nil {
		fmt.Println(er)
		return er
	}
	return nil
}

func (a *ContainerAttacher) holdHijackedConnection(tty bool, inputStream io.Reader, outputStream, errorStream io.Writer, resp types.HijackedResponse) error {
	receiveStdout := make(chan error)
	if outputStream != nil || errorStream != nil {
		go func() {
			receiveStdout <- redirectResponseToOutputStream(tty, outputStream, errorStream, resp.Reader)
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			io.Copy(resp.Conn, inputStream)
		}
		resp.CloseWrite()
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		return err
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			return <-receiveStdout
		}
	}
	return nil
}

func redirectResponseToOutputStream(tty bool, outputStream, errorStream io.Writer, resp io.Reader) error {
	if outputStream == nil {
		outputStream = ioutil.Discard
	}
	if errorStream == nil {
		errorStream = ioutil.Discard
	}
	var err error
	if tty {
		_, err = io.Copy(outputStream, resp)
	} else {
		num, err := stdcopy.StdCopy(outputStream, errorStream, resp)
		log.Println(num, err) // 0 Unrecognized input header: 67
	}
	return err
}
