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
	"github.com/adamzhoul/dockercli/registry"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func proxy2Agent(w http.ResponseWriter, req *http.Request, apiPath string) {

	pathParams := mux.Vars(req)
	cluster := pathParams["cluster"]
	namespace := pathParams["namespace"]
	podName := pathParams["podName"]
	containerName := pathParams["containerName"]
	image, _ := pathParams["image"]
	log.Printf("exec pod: %s, container: %s, namespace: %s, image: %s", podName, containerName, namespace, image)

	// 1. upgrade conn
	pty, err := webterminal.NewTerminalSession(w, req, nil)
	if err != nil {
		ResponseErr(w, err)
		return
	}
	defer func() {
		log.Println("close session.")
		pty.Close()
	}()

	// 2. supply conn params
	var containerImage, containerID, hostIP string
	containerImage, containerID, hostIP, err = registry.Client.FindPodContainerInfo(cluster, namespace, podName, containerName)
	if err != nil {
		pty.Done()
		ResponseErr(w, err)
		return
	}

	//podAgentAddress, err := getAgentAddress(hostIP)
	fmt.Println("get hostIP", hostIP)
	podAgentAddress := "10.245.72.92:18080"
	log.Printf("find pod %s agent address %s", podName, podAgentAddress)
	if err != nil {
		pty.Done()
		ResponseErr(w, err)
		return
	}

	// 3. connect use spdy protocol, link websocket conn and spdy conn
	uri, _ := url.Parse(fmt.Sprintf("http://%s", podAgentAddress))
	uri.Path = apiPath
	params := url.Values{}
	params.Add("attachImage", containerImage)
	params.Add("debugContainerID", containerID)
	uri.RawQuery = params.Encode()
	log.Println("connect to agent ", uri, params)
	exec, err := remotecommand.NewSPDYExecutor(&rest.Config{Host: uri.Host}, "POST", uri)
	if err != nil {
		pty.Done()
		ResponseErr(w, err)
		return
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
	if err != nil {
		pty.Done()
		pty.Close()
		log.Println("stream err:", err)
	}
}

func getAgentAddress(hostIP string) (string, error) {

	agents := kubernetes.FindPodsByLabel(agent.AGENT_NAMESPACE, agent.AGENT_LABEL)
	for _, agent := range agents {

		if agent.Status.HostIP == hostIP {
			return agent.Status.PodIP, nil
		}
	}

	return "", errors.New("agent not found")
}
