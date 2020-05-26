package proxy

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type HTTPConfig struct {
	ListenAddress string
}

type HTTPProxyServer struct {
	server    *http.Server
	config    *HTTPConfig // http server run params
	k8sConfig string      // configs for k8s client-go connection
}

var testAgentAddress string // which agent ip  server proxy to
func NewHTTPProxyServer(config *HTTPConfig, aAddress string) *HTTPProxyServer {

	testAgentAddress = aAddress
	route := proxyRoute()
	return &HTTPProxyServer{
		server: &http.Server{
			Addr:    config.ListenAddress,
			Handler: route},
		config:    config,
		k8sConfig: "",
	}
}

// run and stop
func (s *HTTPProxyServer) Serve(stop chan os.Signal) error {

	go func() {
		log.Printf("Http Server started! Welcome aboard! \n")

		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// stop server
	<-stop

	s.Shutdown()
	return nil
}

func (s *HTTPProxyServer) Shutdown() {

	log.Println("shutting done server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func proxyRoute() *mux.Router {

	//mux := http.NewServeMux()
	route := mux.NewRouter()

	// load html static file
	route.HandleFunc("/", IndexHtml)
	route.HandleFunc("/{action}/ns/{namespace}/pod/{podName}/container/{containerName}", IndexHtml)
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./fe/static/"))))
	//route.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./fe/static/"))))

	// api
	route.HandleFunc("/api/v1/logs", handleLog)
	route.HandleFunc("/api/v1/debug/ns/{namespace}/pod/{podName}/container/{containerName}/image/{image}", handleDebug)
	route.HandleFunc("/api/v1/debug/ns/{namespace}/pod/{podName}/container/{containerName}", handleDebug)
	route.HandleFunc("/api/v1/exec/ns/{namespace}/pod/{podName}/container/{containerName}", handleExec)
	route.HandleFunc("/api/v1/log/ns/{namespace}/pod/{podName}/container/{containerName}", handleLog)
	route.HandleFunc("/healthz", healthz)

	return route
}

// probe health checks
func healthz(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm OK!"))
}

func IndexHtml(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./fe/index.html")
}

func ResponseErr(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	w.Write([]byte(err.Error()))
}
