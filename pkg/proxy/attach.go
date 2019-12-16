package proxy

import (
	"log"
	"net/http"

	"github.com/adamzhoul/dockercli/pkg/kubernetes"
)

// handle websocket connection
// find container node
// connect to agent which on the same node
func handleAttach(w http.ResponseWriter, req *http.Request) {

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
	// exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	// if err != nil {
	// 	return
	// }
	// return exec.Stream(remotecommand.StreamOptions{
	// 	Stdin:             stdin,
	// 	Stdout:            stdout,
	// 	Stderr:            stderr,
	// 	Tty:               tty,
	// 	TerminalSizeQueue: terminalSizeQueue,
	// })
}
