package cmd

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/adamzhoul/dockercli/pkg/agent"
	"github.com/adamzhoul/dockercli/pkg/docker"
	"github.com/spf13/cobra"
)

var (
	listenAddress           string
	dockerAddress           string
	attachTargetContainerID string // for test purpose

	ttyIdleTimeout               time.Duration
	containerGracefulExitTimeout time.Duration
	containerCreateTimeout       time.Duration
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

	var ttyIdleTimeoutMinute, containerCreateTimeoutSecond, containerGracefulExitTimeoutSecond int
	agentCmd.Flags().IntVar(&ttyIdleTimeoutMinute, "ttyTimeout", 30, "tty connect idle timeout in minute")
	agentCmd.Flags().IntVar(&containerCreateTimeoutSecond, "containerCreateTimeout", 15, "container create timeout in second")
	agentCmd.Flags().IntVar(&containerGracefulExitTimeoutSecond, "containerExitTimeout", 15, "container exit timeout in second")

	ttyIdleTimeout = time.Duration(ttyIdleTimeoutMinute) * time.Minute
	containerCreateTimeout = time.Duration(containerCreateTimeoutSecond) * time.Second
	containerGracefulExitTimeout = time.Duration(containerGracefulExitTimeoutSecond) * time.Second
}

func runAgent(cmd *cobra.Command, args []string) error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)


	runtimeConfig := docker.RuntimeConfig{
		DockerEndpoint:        dockerAddress,
		StreamIdleTimeout:     ttyIdleTimeout,
		StreamCreationTimeout: containerCreateTimeout,
		GracefulExitTimeout:   containerGracefulExitTimeout,
	}
	log.Println("load config success:", listenAddress, dockerAddress, attachTargetContainerID, runtimeConfig)
	docker.InitDockerclientConn(runtimeConfig)

	// start an HttpServer
	agentServer := agent.NewHTTPAgentServer(listenAddress, runtimeConfig, attachTargetContainerID)
	agentServer.Serve(stop)

	return nil
}

func agentExec() {

	err := agentCmd.Execute()
	if err != nil {
		log.Println(err)
	}

	// quit, close conns
	docker.CloseConn()
}
