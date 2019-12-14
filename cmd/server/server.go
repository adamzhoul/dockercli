package main

import (
	web "github.com/adamzhoul/dockercli/web"
	"os"
	"os/signal"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start an HttpServer
	web.RunHttpServer(stop)
	// docker.imag()
}
