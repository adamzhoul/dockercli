package agent

import (
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

func NewHTTPAgentServer(config *HTTPConfig) *HTTPAgentServer {

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
		log.Printf("Http Server started! Welcome aboard! \n")

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
