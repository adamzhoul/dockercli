package proxy

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/adamzhoul/dockercli/pkg/agent"
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
	if err != nil {
		ResponseErr(w, err)
		return
	}

	if agentAddress == "" {
		agentAddress, err = getAgentAddress("mservice", "96143-helloworld-mservice-557545669f-drqdf")
		if err != nil {
			ResponseErr(w, err)
			return
		}
	}
	if agentAddress == "" {
		ResponseErr(w, errors.New("can't find agent ip"))
		return
	}

	// connect use spdy protocol, link websocket conn and spdy conn
	uri, err := url.Parse(fmt.Sprintf("http://%s", agentAddress))
	if err != nil {
		return
	}
	uri.Path = fmt.Sprintf("/api/v1/debug")
	config := rest.Config{Host: uri.Host}
	exec, err := remotecommand.NewSPDYExecutor(&config, "POST", uri)
	if err != nil {
		log.Println(err)
		return
	}

	exec.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
}

func getAgentAddress(namespace string, podName string) (string, error) {

	// 1. find container node
	pod := kubernetes.FindPodByName(namespace, podName)
	if pod == nil {
		return "", errors.New("pod not found")
	}

	// 2. get agent ip
	agents := kubernetes.FindPodsByLabel(agent.AGENT_NAMESPACE, agent.AGENT_LABEL)
	for _, agent := range agents {

		if agent.Status.HostIP == pod.Status.HostIP {
			return agent.Status.PodIP, nil
		}
	}

	return "", errors.New("agent not found")
}
