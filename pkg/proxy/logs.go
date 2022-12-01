package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	util "github.com/adamzhoul/dockercli/pkg"
	"github.com/adamzhoul/dockercli/registry"
	"github.com/gorilla/mux"
)

func handleLog(w http.ResponseWriter, req *http.Request) {
	reverseProxy2Agent(w, req, "/api/v1/log")
}

func reverseProxy2Agent(w http.ResponseWriter, req *http.Request, apiPath string) {

	pathParams := mux.Vars(req)
	cluster := pathParams["cluster"]
	namespace := pathParams["namespace"]
	podName := pathParams["podName"]
	containerName := pathParams["containerName"]
	logger, ok := req.Context().Value("logger").(util.ShellLogger)
	if !ok {
		ResponseErr(w, errors.New("log error"))
		return
	}
	logger.Info(fmt.Sprintf("exec pod: %s, container: %s, namespace: %s, image: %s", podName, containerName, namespace, "image"))

	// 2. supply conn params
	var containerID, hostIP string
	_, containerID, hostIP, err := registry.Client.FindPodContainerInfo(cluster, namespace, podName, containerName)
	if err != nil {
		logger.Info(err.Error())
		ResponseErr(w, err)
		return
	}

	logger.Info("get hostIP", hostIP)
	//podAgentAddress, err := registry.Client.FindAgentIp(cluster, hostIP)
	podAgentAddress := hostIP
	if err != nil {
		ResponseErr(w, err)
		return
	}

	// 3. connect use spdy protocol, link websocket conn and spdy conn
	uri, _ := url.Parse(fmt.Sprintf("http://%s:%s", podAgentAddress, registry.Client.FindAgentPort()))
	//uri.RawPath = apiPath
	params := url.Values{}
	params.Add("debugContainerID", containerID)
	uri.RawQuery = params.Encode()

	reverseProxy(w, req, uri, apiPath)
	logger.Info("reverseProxy log done")
}

func reverseProxy(w http.ResponseWriter, req *http.Request, uri *url.URL, apiPath string) {
	req.URL.Path = apiPath
	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.ServeHTTP(w, req)
}
