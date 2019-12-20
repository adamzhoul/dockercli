package agent

import (
	"github.com/adamzhoul/dockercli/pkg/docker"
	"net/http"
)

func (s *HTTPAgentServer) handleLog(w http.ResponseWriter, req *http.Request) {

	debugContainerID := req.FormValue("debugContainerID")

	//client.ContainerLogs(context.Background(), debugContainerID)
	docker.TailLog(debugContainerID)

}
