package proxy

import (
	"context"
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

func NewHTTPProxyServer(config *HTTPConfig) *HTTPProxyServer {

	muex := proxyRoute()
	return &HTTPProxyServer{
		server: &http.Server{
			Addr:    config.ListenAddress,
			Handler: muex},
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

	Test()

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

func Err(w http.ResponseWriter, err string) {
	w.Write([]byte(err))
}

func proxyRoute() *http.ServeMux {

	mux := http.NewServeMux()

	// load html static file
	mux.HandleFunc("/", IndexHtml)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./fe/static/"))))

	// api
	mux.HandleFunc("/api/v1/logs", handleLog)
	mux.HandleFunc("/api/v1/attach", handleAttach)
	mux.HandleFunc("/healthz", healthz)

	return mux
}

// probe health checks
func healthz(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm OK!"))
}

func IndexHtml(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./fe/index.html")
}
