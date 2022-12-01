package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func (s *HTTPAgentServer) handleFile(w http.ResponseWriter, req *http.Request) {

	debugContainerID := req.FormValue("debugContainerID")
	debugContainerID = strings.TrimLeft(debugContainerID, "containerd://")
	fmt.Println("file exec into container:", debugContainerID)

	file := req.FormValue("file")
	cmd := []string{"tar", "cf", "-", file}

	// tar only support relavite path
	if filepath.IsAbs(file) { // using -C DIR Change to DIR before operation
		dir, fileName := filepath.Split(file)
		cmd = []string{"tar", "cf", "-", "-C", dir, fileName}
	}

	fmt.Println("copy command: ", cmd)
	rp, err := s.runtimeService.Exec(&runtimeapi.ExecRequest{
		Stdin:       true, // tell when conn is closed
		Stdout:      true,
		Stderr:      true,
		Tty:         false,
		ContainerId: debugContainerID,
		Cmd:         cmd,
	})
	if err != nil {
		fmt.Println("exec err:", err)
		return
	}

	u, _ := url.Parse(rp.Url)
	proxyStream(w, req, u)
	fmt.Println("file copy done")
}
