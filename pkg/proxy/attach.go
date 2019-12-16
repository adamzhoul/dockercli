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

	podsName := req.PostForm["podsName"]

	if len(podsName) != 1 {
		Err(w, "only one pod supported!")
		return
	}

	// 1. find container node
	pods := kubernetes.FindPodsByName("", podsName)
	if len(pods) != 1 {
		return
	}

	// 2. get agent ip , connect use spdy protocol
	agent := kubernetes.FindPodsByLabel("", "")
	if len(agent) != 1 {

	}
	ip := agent[0].Status.PodIP
	log.Println("connect to agetn :", ip)

	// 3. link websocket conn and spdy conn

}
