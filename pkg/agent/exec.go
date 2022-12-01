package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"k8s.io/apimachinery/pkg/util/proxy"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog/v2"
)

func (s *HTTPAgentServer) handleCriExec(w http.ResponseWriter, req *http.Request) {

	debugContainerID := req.FormValue("debugContainerID")
	debugContainerID = strings.TrimLeft(debugContainerID, "containerd://")
	fmt.Println("exec into container:", debugContainerID)

	rp, err := s.runtimeService.Exec(&runtimeapi.ExecRequest{
		Stdin:       true,
		Stdout:      true,
		Tty:         true,
		ContainerId: debugContainerID,
		Cmd:         []string{"/bin/sh"},
	})
	if err != nil {
		fmt.Println("exec err:", err)
		return
	}

	u, _ := url.Parse(rp.Url)

	proxyStream(w, req, u)
}

type responder struct{}

func (r *responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.ErrorS(err, "Error while proxying request")
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// proxyStream proxies stream to url.
func proxyStream(w http.ResponseWriter, r *http.Request, url *url.URL) {
	// TODO(random-liu): Set MaxBytesPerSec to throttle the stream.
	handler := proxy.NewUpgradeAwareHandler(url, nil /*transport*/, false /*wrapTransport*/, true /*upgradeRequired*/, &responder{})
	handler.ServeHTTP(w, r)
}
