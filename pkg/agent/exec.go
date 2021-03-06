package agent

import (
	"net/http"

	"github.com/adamzhoul/dockercli/pkg/docker"
	remoteapi "k8s.io/apimachinery/pkg/util/remotecommand"
	kubeletremote "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
)

func (s *HTTPAgentServer) handleExec(w http.ResponseWriter, req *http.Request) {

	debugContainerID := req.FormValue("debugContainerID")
	if !auth(req){
		http.Error(w, "Unauthorized", 401)
		return
	}

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
		GetExecAttacher(),
		"",
		"",
		debugContainerID,
		streamOpts,
		s.RuntimeConfig.StreamIdleTimeout, // idle timeout will lead server send fin package 
		s.RuntimeConfig.StreamCreationTimeout,
		remoteapi.SupportedStreamingProtocols)
}

func GetExecAttacher() *docker.ContainerExecAttacher {

	return docker.NewContainerExecAttacher()
}
