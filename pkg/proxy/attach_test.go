package proxy

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func TestSpdy(t *testing.T) {
	uri, err := url.Parse(fmt.Sprintf("http://127.0.0.1:8090"))
	if err != nil {
		return
	}
	uri.Path = fmt.Sprintf("/api/v1/attach")
	config := rest.Config{Host: fmt.Sprintf("http://127.0.0.1:8090")}
	exec, err := remotecommand.NewSPDYExecutor(&config, "POST", uri)
	if err != nil {
		return
	}
	exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
		//TerminalSizeQueue: terminalSizeQueue,
	})
}
