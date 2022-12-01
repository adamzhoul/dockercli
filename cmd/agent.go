package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/pkg/agent"
	"github.com/spf13/cobra"
)

var (
	listenAddress           string
	dockerAddress           string
	attachTargetContainerID string // for test purpose
)

var agentCmd = &cobra.Command{
	Use:           "agent",
	Short:         "agent is a command line tool which connect to docker",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runAgent,
}

func init() {

	agentCmd.Flags().StringVar(&listenAddress, "listenAddress", "0.0.0.0:80", "http listener")
	agentCmd.Flags().StringVar(&dockerAddress, "dockerAddress", "/var/run/docker.sock", "docker socket path")
	agentCmd.Flags().StringVar(&attachTargetContainerID, "cid", "", "which container attach to")

}

func runAgent(cmd *cobra.Command, args []string) error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Println("load config success:", listenAddress, dockerAddress, attachTargetContainerID)

	// start an HttpServer
	agentServer := agent.NewHTTPAgentServer(listenAddress, attachTargetContainerID)
	agentServer.Serve(stop)

	return nil
}
