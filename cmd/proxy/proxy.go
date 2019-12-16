package main

import (
	"github.com/adamzhoul/dockercli/pkg/proxy"
	"os"
	"os/signal"
)

func initConfig() {

}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	config := proxy.HTTPConfig{
		ListenAddress: "0.0.0.0:8089",
	}
	// start an HttpServer
	proxy := proxy.NewHTTPProxyServer(&config)
	proxy.Serve(stop)

	// docker.imag()
}
