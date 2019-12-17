package main

import (
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/pkg/agent"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	config := agent.HTTPConfig{
		ListenAddress: "0.0.0.0:8090",
	}

	// start an HttpServer
	agentServer := agent.NewHTTPAgentServer(&config)
	agentServer.Serve(stop)
}
