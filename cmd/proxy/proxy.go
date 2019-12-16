package main

import (
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/pkg/kubernetes"
	"github.com/adamzhoul/dockercli/pkg/proxy"
)

func initConfig() {

	kubernetes.InitClientgo("./configs/kube/config")
}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	initConfig()

	config := proxy.HTTPConfig{
		ListenAddress: "0.0.0.0:8089",
	}

	// start an HttpServer
	proxy := proxy.NewHTTPProxyServer(&config)
	proxy.Serve(stop)

	// docker.imag()
}
