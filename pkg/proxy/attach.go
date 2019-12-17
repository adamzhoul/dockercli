package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/adamzhoul/dockercli/pkg/kubernetes"
	"github.com/adamzhoul/dockercli/pkg/webterminal"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// handle websocket connection
// find container node
// connect to agent which on the same node
func handleAttach(w http.ResponseWriter, req *http.Request) {
	pty, err := webterminal.NewTerminalSession(w, req, nil)

	// 1. upgrade protocol to websocket
	podsName := req.PostForm["podsName"]

	if len(podsName) != 1 {
		Err(w, "only one pod supported!")
		return
	}

	// 2. find container node
	pods := kubernetes.FindPodsByName("", podsName)
	if len(pods) != 1 {
		return
	}

	// 3. get agent ip
	agent := kubernetes.FindPodsByLabel("", "")
	if len(agent) != 1 {

	}
	ip := agent[0].Status.PodIP
	log.Println("connect to agetn :", ip)

	// 4. connect use spdy protocol, link websocket conn and spdy conn
	uri, err := url.Parse(fmt.Sprintf("http://%s:%d", ip, 8090))
	if err != nil {
		return
	}
	uri.Path = fmt.Sprintf("/api/v1/debug")
	config := rest.Config{Host: fmt.Sprintf("http://%s:%d", ip, 8090)}
	exec, err := remotecommand.NewSPDYExecutor(&config, "POST", uri)
	if err != nil {
		return
	}
	exec.Stream(remotecommand.StreamOptions{
		Stdin:  pty,
		Stdout: pty,
		Stderr: pty,
		Tty:    true,
		//TerminalSizeQueue: terminalSizeQueue,
	})
}

func Test() {
	uri, err := url.Parse(fmt.Sprintf("http://127.0.0.1:8090"))
	if err != nil {
		return
	}
	uri.Path = fmt.Sprintf("/api/v1/debug")
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
