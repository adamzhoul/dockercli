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
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// handle websocket connection
// find container node
// connect to agent which on the same node
func handleDebug(w http.ResponseWriter, req *http.Request) {

	pathParams := mux.Vars(req)
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

	// 2. supply conn params
	var containerImage, containerID, hostIP string
	var podAgentAddress string
	if testAgentAddress == "" {
		containerImage, containerID, hostIP, err = findPodContainerInfo(namespace, podName, containerName)
		if err != nil {
			ResponseErr(w, err)
			return
		}

		podAgentAddress, err = getAgentAddress(hostIP)
		log.Printf("find pod %s agent address %s", podName, podAgentAddress)
		if err != nil {
			ResponseErr(w, err)
			return
		}
	} else {
		podAgentAddress = testAgentAddress
	}

	// 3. connect use spdy protocol, link websocket conn and spdy conn
	uri, _ := url.Parse(fmt.Sprintf("http://%s", podAgentAddress))
	uri.Path = fmt.Sprintf("/api/v1/debug")
	params := url.Values{}
	params.Add("attachImage", containerImage)
	params.Add("debugContainerID", containerID)
	uri.RawQuery = params.Encode()
	// config := rest.Config{Host: uri.Host}
	log.Println("connect to agent ", uri, params)
	exec, err := remotecommand.NewSPDYExecutor(&rest.Config{Host: uri.Host}, "POST", uri)
	if err != nil {
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
		log.Println(err)
	}
}

// get pod container info
// include: containerImage containerID HostIP
func findPodContainerInfo(namespace string, podName string, containerName string) (string, string, string, error) {

	var image, containerID string

	// 1. find pod
	pod := kubernetes.FindPodByName(namespace, podName)
	if pod == nil {
		return "", "", "", errors.New("pod not found")
	}

	// 2. find container image
	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			image = container.Image
			break
		}
	}

	// 3. find container ID
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == containerName {
			containerID = containerStatus.ContainerID
			break
		}
	}

	if image == "" || containerID == "" {
		return image, containerID, pod.Status.HostIP, errors.New("pod info error ")
	}

	return image, containerID, pod.Status.HostIP, nil

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
