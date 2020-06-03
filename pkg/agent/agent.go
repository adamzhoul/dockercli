package agent

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adamzhoul/dockercli/pkg/docker"
)

type HTTPAgentServer struct {
	server        *http.Server
	ListenAddress string
	RuntimeConfig docker.RuntimeConfig
}

// for test purpose
var testAttachTargetContainerID string

const (
	AGENT_NAMESPACE = "default"
	AGENT_LABEL     = "component=dockercli.agent"
)

func NewHTTPAgentServer(addr string, runtimeConfig docker.RuntimeConfig, attachTargetContainerID string) *HTTPAgentServer {

	testAttachTargetContainerID = attachTargetContainerID
	s := &HTTPAgentServer{
		server: &http.Server{
			Addr: addr,
		},
		RuntimeConfig: runtimeConfig,
	}
	s.server.Handler = s.proxyRoute()

	return s
}

// run and stop
func (s *HTTPAgentServer) Serve(stop chan os.Signal) error {

	go func() {
		log.Printf(fmt.Sprintf("Http Server started at %s! Welcome aboard! \n", s.ListenAddress))

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

	mux.HandleFunc("/api/v1/debug", s.handleDebug)
	mux.HandleFunc("/api/v1/exec", s.handleExec)
	mux.HandleFunc("/api/v1/log", s.handleLog)

	return mux
}

// probe health checks
func (s *HTTPAgentServer) healthz(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm OK!"))
}

func (s *HTTPAgentServer) Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Welcome!"))
}

func ResponseErr(w http.ResponseWriter, err error, code int) {
	log.Println(err.Error())
	http.Error(w, err.Error(), code)
}
