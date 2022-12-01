package proxy

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/adamzhoul/dockercli/pkg/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func TestGetAgentAddress(t *testing.T) {
	kubernetes.InitClientgo("../../configs/kube/config")
	agentAddress, err := getAgentAddress("")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(agentAddress)
}

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

func TestUrl(t *testing.T) {
	uri, _ := url.Parse(fmt.Sprintf("http://%s:%s", "127.0.0.1", "8080"))
	uri.Path = "/api/v1/log"

	fmt.Printf("get uri: %+v", uri)
}
