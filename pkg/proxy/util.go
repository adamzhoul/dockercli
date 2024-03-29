package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	util "github.com/adamzhoul/dockercli/pkg"
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
	// image, _ := pathParams["image"]
	logger, ok := req.Context().Value("logger").(util.ShellLogger)
	if !ok {
		ResponseErr(w, errors.New("log error"))
		return
	}
	logger.Info(fmt.Sprintf("exec pod: %s, container: %s, namespace: %s, image: %s", podName, containerName, namespace, "image"))

	// 1. upgrade conn
	pty, err := webterminal.NewTerminalSession(w, req, nil)
	if err != nil {
		ResponseErr(w, err)
		return
	}
	defer func() {
		logger.Info("close session.")
		pty.Close()
	}()

	// 2. supply conn params
	var containerImage, containerID, hostIP string
	containerImage, containerID, hostIP, err = registry.Client.FindPodContainerInfo(cluster, namespace, podName, containerName)
	if err != nil {
		logger.Info(err.Error())
		pty.Done()
		ResponseErr(w, err)
		return
	}

	logger.Info("get hostIP", hostIP)
	//podAgentAddress, err := registry.Client.FindAgentIp(cluster, hostIP)
	podAgentAddress := hostIP
	if err != nil {
		pty.Done()
		ResponseErr(w, err)
		return
	}

	// 3. connect use spdy protocol, link websocket conn and spdy conn
	uri, _ := url.Parse(fmt.Sprintf("http://%s:%s", podAgentAddress, registry.Client.FindAgentPort()))
	uri.Path = apiPath
	params := url.Values{}
	params.Add("attachImage", containerImage)
	params.Add("debugContainerID", containerID)
	uri.RawQuery = params.Encode()

	username := req.Context().Value("username")
	password := util.EncryptionArithmetic(username.(string), "oasdf923n")
	logger.Info("request verify:", username.(string), ",", password)

	exec, err := remotecommand.NewSPDYExecutor(&rest.Config{Host: uri.Host, Username: username.(string), Password: password}, "POST", uri)
	if err != nil {
		logger.Info(err.Error())
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
	if err != nil { // 这里断开了，会往标准输出写么？？？？
		pty.Done()
		pty.Close()
		logger.Info("stream err:", err.Error())
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
