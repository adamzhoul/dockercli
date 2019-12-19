package docker

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"time"

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

// attach to container
// 1. handle size
// 2. attach
// 3. hold conn
func (a *ContainerAttacher) AttachContainer(name string, uid kubetype.UID, container string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize) error {

	// handle size
	HandleResizing(resize, func(size remotecommand.TerminalSize) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		er := a.client.ContainerResize(ctx, container, types.ResizeOptions{
			Height: uint(size.Height),
			Width:  uint(size.Width),
		})
		if er != nil {
			log.Println("resize failed:", er)
		}

	})

	// attach to container
	opts := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stderr: true,
		Stdout: true,
	}
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	log.Println("attach container:", container)
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
		log.Println(er)
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

func HandleResizing(resize <-chan remotecommand.TerminalSize, resizeFunc func(size remotecommand.TerminalSize)) {
	if resize == nil {
		return
	}

	go func() {
		//defer runtime.HandleCrash()

		for size := range resize {
			if size.Height < 1 || size.Width < 1 {
				continue
			}
			resizeFunc(size)
		}
	}()
}
