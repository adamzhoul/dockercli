package agent

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type HTTPConfig struct {
	ListenAddress string
}

type HTTPAgentServer struct {
	server *http.Server
	config *HTTPConfig // http server run params
}

// for test purpose
var attachDebugTargetContainerID string

func NewHTTPAgentServer(config *HTTPConfig, attachTargetContainerID string) *HTTPAgentServer {

	attachDebugTargetContainerID = attachTargetContainerID
	muex := proxyRoute()
	return &HTTPAgentServer{
		server: &http.Server{
			Addr:    config.ListenAddress,
			Handler: muex},
		config: config,
	}
}

// run and stop
func (s *HTTPAgentServer) Serve(stop chan os.Signal) error {

	go func() {
		log.Printf(fmt.Sprintf("Http Server started at %s! Welcome aboard! \n", s.config.ListenAddress))

		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// stop server
	<-stop

	//s.Shutdown()
	return nil
}

func proxyRoute() *http.ServeMux {

	mux := http.NewServeMux()

	// load welcome page
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/healthz", healthz)

	mux.HandleFunc("/api/v1/debug", handleDebug)

	return mux
}

// probe health checks
func healthz(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm OK!"))
}

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Welcome!"))
}
