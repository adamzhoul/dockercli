package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adamzhoul/dockercli/pkg/auth"
	"github.com/gorilla/mux"
)

type HTTPConfig struct {
	ListenAddress string
}

type HTTPProxyServer struct {
	server *http.Server
	config *HTTPConfig // http server run params
}

type httpHandler struct {
	r *mux.Router
}

// extract from url is really really not a good design
// todo: fix it
func extractResourceActionFromUrl(req *http.Request) (resource string, action string) {
	rawUrl := req.URL.Path
	cluster := ""
	namespace := ""
	pod := ""
	items := strings.Split(rawUrl, "/")

	if !strings.HasPrefix(rawUrl, "/api") { // html page, /{action}/cluster/{cluster}/ns/{namespace}/pod/{podName}/c....
		fmt.Println(items)
		action = items[1]
		cluster = items[3]
		namespace = items[5]
		pod = items[7]
	} else {
		action = items[3]
		cluster = items[5]
		namespace = items[7]
		pod = items[9]
	}

	arr := strings.Split(pod, "-")
	pod = strings.Join(arr[:len(arr)-2], "-")

	resource = fmt.Sprintf("/%s/%s/%s", cluster, namespace, pod)
	return
}

// /static 、 / 、 /healthz
func passTokenCheck(req *http.Request) bool {

	if strings.HasPrefix(req.URL.Path, "/static") ||
		req.URL.RawPath == "/" || req.URL.RawPath == "/healthz" {
		return true
	}

	return false
}

func (h *httpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// pass allowList: /static 、 / 、 /healthz
	if passTokenCheck(req) {
		h.r.ServeHTTP(rw, req)
		return
	}

	// get userinfo before action
	token, err := req.Cookie("token")
	if token == nil || err != nil {
		//http.Redirect(rw, req, fmt.Sprintf("http://www.google.com?redirect=%s", req.URL.RawPath), 302)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("401 HTTP status code returned!"))
		return
	}

	// check auth before action
	resource, action := extractResourceActionFromUrl(req)
	_, pass := auth.CheckUser(token.Value, resource, action)
	if !pass {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("403 HTTP status code returned!"))
		return
	}

	h.r.ServeHTTP(rw, req)
}

func NewHTTPProxyServer(config *HTTPConfig) *HTTPProxyServer {

	routeHandler := proxyRoute()
	return &HTTPProxyServer{
		server: &http.Server{
			Addr:    config.ListenAddress,
			Handler: routeHandler},
		config: config,
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

func proxyRoute() *httpHandler {

	//mux := http.NewServeMux()
	route := mux.NewRouter()

	// load html static file
	route.HandleFunc("/", IndexHtml)
	route.HandleFunc("/{action}/cluster/{cluster}/ns/{namespace}/pod/{podName}/container/{containerName}", IndexHtml)
	route.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./fe/static/"))))

	// api
	route.HandleFunc("/api/v1/debug/cluster/{cluster}/ns/{namespace}/pod/{podName}/container/{containerName}/image/{image}", handleDebug)
	route.HandleFunc("/api/v1/debug/cluster/{cluster}/ns/{namespace}/pod/{podName}/container/{containerName}", handleDebug)
	route.HandleFunc("/api/v1/exec/cluster/{cluster}/ns/{namespace}/pod/{podName}/container/{containerName}", handleExec)
	route.HandleFunc("/api/v1/log/cluster/{cluster}/ns/{namespace}/pod/{podName}/container/{containerName}", handleLog)
	route.HandleFunc("/healthz", healthz)

	return &httpHandler{
		r: route,
	}
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
