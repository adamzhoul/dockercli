package cmd

import (
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/pkg/kubernetes"
	"github.com/adamzhoul/dockercli/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	agentAddress       string
	proxyListenAddress string
	kubeConfigPath     string
	skipKube           bool
)

var proxyCmd = &cobra.Command{
	Use:           "proxy",
	Short:         "proxy is a command line tool which connect to docker",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runProxy,
}

func init() {

	proxyCmd.Flags().StringVar(&agentAddress, "agent", "", "agent ip port")
	proxyCmd.Flags().StringVar(&proxyListenAddress, "addr", "0.0.0.0:80", "http listener")
	proxyCmd.Flags().StringVar(&kubeConfigPath, "kubeConfig", "./configs/kube/config", "kube config ")
	proxyCmd.Flags().BoolVar(&skipKube, "skipKube", false, "skip kube config or not")

}

func runProxy(cmd *cobra.Command, args []string) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	initConfig()

	config := proxy.HTTPConfig{
		ListenAddress: proxyListenAddress,
	}

	// start an HttpServer
	proxy := proxy.NewHTTPProxyServer(&config, agentAddress)
	proxy.Serve(stop)

	// docker.imag()
	return nil
}

func initConfig() {

	if !skipKube {
		kubernetes.InitClientgo(kubeConfigPath)
	}

}
