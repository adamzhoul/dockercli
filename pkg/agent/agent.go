package agent

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	util "github.com/adamzhoul/dockercli/pkg"
	cri "k8s.io/cri-api/pkg/apis"
	"k8s.io/kubernetes/pkg/kubelet/cri/remote"
)

type HTTPAgentServer struct {
	server        *http.Server
	ListenAddress string
	// RuntimeConfig runtime.RuntimeConfig

	runtimeService cri.RuntimeService
}

// for test purpose
// var testAttachTargetContainerID string

const (
	AGENT_NAMESPACE = "default"
	AGENT_LABEL     = "component=dockercli.agent"
)

func NewHTTPAgentServer(addr string, attachTargetContainerID string) *HTTPAgentServer {

	runtimeService, err := remote.NewRemoteRuntimeService("unix:///var/run/containerd/containerd.sock", time.Minute)
	if err != nil {
		fmt.Println("new runtime err:", err)
		return nil
	}

	s := &HTTPAgentServer{
		server: &http.Server{
			Addr: addr,
		},
		runtimeService: runtimeService,
	}
	s.server.Handler = s.proxyRoute()

	return s
}

// run and stop
func (s *HTTPAgentServer) Serve(stop chan os.Signal) error {

	go func() {
		log.Printf("Http Server started at %s! Welcome aboard! \n", s.ListenAddress)

		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// stop server
	<-stop

	//s.Shutdown()
	return nil
}

func (s *HTTPAgentServer) proxyRoute() *http.ServeMux {

	mux := http.NewServeMux()

	// load welcome page
	mux.HandleFunc("/", s.Index)
	mux.HandleFunc("/healthz", s.healthz)

	// mux.HandleFunc("/api/v1/debug", s.handleDebug)
	mux.HandleFunc("/api/v1/exec", s.handleCriExec)
	mux.HandleFunc("/api/v1/log", s.handleCriLog)
	mux.HandleFunc("/api/v1/file", s.handleFile)

	return mux
}

// probe health checks
func (s *HTTPAgentServer) healthz(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm OK!"))
}

func (s *HTTPAgentServer) Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("welcome ", req.URL)
	w.Write([]byte("Welcome!"))
}

func ResponseErr(w http.ResponseWriter, err error, code int) {
	log.Println(err.Error())
	http.Error(w, err.Error(), code)
}

func auth(req *http.Request) bool {

	username, password, ok := req.BasicAuth()
	if !ok {
		return false
	}
	log.Println("request user:", username, password)
	supposedPassword := util.EncryptionArithmetic(username, "oasdf923n")

	return supposedPassword == password
}
