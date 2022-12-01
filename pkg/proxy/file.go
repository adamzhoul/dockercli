package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	util "github.com/adamzhoul/dockercli/pkg"
	"github.com/adamzhoul/dockercli/registry"
	"github.com/gorilla/mux"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func handleFile(w http.ResponseWriter, req *http.Request) {
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
	var containerImage, containerID, hostIP string
	containerImage, containerID, hostIP, err := registry.Client.FindPodContainerInfo(cluster, namespace, podName, containerName)
	if err != nil {
		ResponseErr(w, err)
		return
	}

	logger.Info("get hostIP", hostIP, containerImage)
	//podAgentAddress, err := registry.Client.FindAgentIp(cluster, hostIP)
	podAgentAddress := hostIP
	if err != nil {
		ResponseErr(w, err)
		return
	}

	uri, _ := url.Parse(fmt.Sprintf("http://%s:%s", podAgentAddress, registry.Client.FindAgentPort()))
	uri.Path = "/api/v1/file"
	params := url.Values{}
	params.Add("debugContainerID", containerID)
	params.Add("file", req.URL.Query().Get("file"))
	uri.RawQuery = params.Encode()

	username := req.Context().Value("username")
	password := util.EncryptionArithmetic(username.(string), "oasdf923n")
	logger.Info("request verify:", username.(string), ",", password)

	exec, err := remotecommand.NewSPDYExecutor(&rest.Config{Host: uri.Host, Username: username.(string), Password: password}, "POST", uri)
	if err != nil {
		logger.Info(err.Error())
		ResponseErr(w, err)
		return
	}

	wr := wrpRequest{
		R: req,
	}
	errW := httptest.NewRecorder()

	w.Header().Add("Content-type", "application/x-tar")
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  wr,
		Stdout: w,
		Stderr: errW,
		Tty:    false,
	})
	if err != nil {
		logger.Info("stream file err:", err)
	}

	errMsg := []byte{}
	n, err := errW.Body.Read(errMsg)
	logger.Info("read stderr ", n, err, string(errMsg))

	// read from w
	logger.Info("proxy file done")
}

type wrpRequest struct {
	R *http.Request
}

// todo when tar takes time and client go away
// make stdout, stderr err at same time to exit goroutine in exec.Stream
/*
	var wg sync.WaitGroup
	p.copyStdout(&wg)
	p.copyStderr(&wg)
	// we're waiting for stdout/stderr to finish copying
	wg.Wait()
*/
func (wr wrpRequest) Read(p []byte) (int, error) {
	<-wr.R.Context().Done()
	fmt.Println("req ctx is done")

	// todo we should close response to exit copy in exec.stream

	return 0, fmt.Errorf("client finished")
}
