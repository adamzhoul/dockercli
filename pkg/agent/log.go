package agent

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/util/flushwriter"
	"k8s.io/kubernetes/pkg/kubelet/kuberuntime/logs"
)

func (s *HTTPAgentServer) handleCriLog(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	debugContainerID := req.FormValue("debugContainerID")
	debugContainerID = strings.TrimLeft(debugContainerID, "containerd://")
	fmt.Println("log container:", debugContainerID)

	status, err := s.runtimeService.ContainerStatus(debugContainerID)
	if err != nil {
		fmt.Println("get container status err:", err)
		return
	}

	TailLines := int64(200)
	opts := logs.NewLogOptions(&corev1.PodLogOptions{
		Follow:    true,
		TailLines: &TailLines,
	}, time.Now())
	fw := flushwriter.Wrap(w)
	w.Header().Set("Transfer-Encoding", "chunked")
	err = logs.ReadLogs(ctx, status.GetLogPath(), debugContainerID, opts, s.runtimeService, fw, fw)

	fmt.Println("read log done ", err)
}
