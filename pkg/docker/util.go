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
	"k8s.io/client-go/tools/remotecommand"
)

func holdHijackedConnection(tty bool, inputStream io.Reader, outputStream, errorStream io.Writer, resp types.HijackedResponse) error {

	receiveStdout := make(chan error)
	if outputStream != nil || errorStream != nil {
		go func() {
			receiveStdout <- redirectResponseToOutputStream(tty, outputStream, errorStream, resp.Reader)
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			n, err := io.Copy(resp.Conn, inputStream)
			log.Println("input  number ", n, err)
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
		n, err := io.Copy(outputStream, resp)
		log.Println("output  number ", n, err)
	} else {
		num, err := stdcopy.StdCopy(outputStream, errorStream, resp)
		log.Println(num, err) // 0 Unrecognized input header: 67
	}
	return err
}

func handleResizing(resize <-chan remotecommand.TerminalSize, client *dockerclient.Client, container string, resizeFunc func(size remotecommand.TerminalSize, client *dockerclient.Client, container string)) {
	if resize == nil {
		return
	}

	go func() {
		//defer runtime.HandleCrash()

		for size := range resize {
			if size.Height < 1 || size.Width < 1 {
				continue
			}
			resizeFunc(size, client, container)
		}
	}()
}

func resizeContainer(size remotecommand.TerminalSize, client *dockerclient.Client, container string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	er := client.ContainerResize(ctx, container, types.ResizeOptions{
		Height: uint(size.Height),
		Width:  uint(size.Width),
	})
	if er != nil {
		log.Println("resize failed:", er)
	}

}

func resizeExecContainer(size remotecommand.TerminalSize, client *dockerclient.Client, container string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	er := client.ContainerExecResize(ctx, container, types.ResizeOptions{
		Height: uint(size.Height),
		Width:  uint(size.Width),
	})
	if er != nil {
		log.Println("resize failed:", er)
	}

}
