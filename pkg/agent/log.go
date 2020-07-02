package agent

import (
	"net/http"

	"github.com/adamzhoul/dockercli/pkg/docker"
	remoteapi "k8s.io/apimachinery/pkg/util/remotecommand"
	kubeletremote "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
)

func (s *HTTPAgentServer) handleLog(w http.ResponseWriter, req *http.Request) {

	if !auth(req){
		http.Error(w, "Unauthorized", 401)
		return
	}

	debugContainerID := req.FormValue("debugContainerID")

	// 2. attach to container
	streamOpts := &kubeletremote.Options{
		Stdin:  true,
		Stdout: true,
		Stderr: false,
		TTY:    true,
	}
	kubeletremote.ServeAttach(
		w,
		req,
		GetLogAttacher(),
		"",
		"",
		debugContainerID,
		streamOpts,
		s.RuntimeConfig.StreamIdleTimeout, // idle timeout will lead server send fin package 
		s.RuntimeConfig.StreamCreationTimeout,
		remoteapi.SupportedStreamingProtocols)
}

func GetLogAttacher() *docker.ContainerLogAttacher {

	return docker.NewContainerLogAttacher()
}
