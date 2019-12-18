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
)

var proxyCmd = &cobra.Command{
	Use:           "proxy",
	Short:         "agent is a command line tool which connect to docker",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runProxy,
}

func init() {

	proxyCmd.Flags().StringVar(&agentAddress, "agentAddress", "", "agent ip port")
	proxyCmd.Flags().StringVar(&proxyListenAddress, "proxyListenAddress", "0.0.0.0:80", "http listener")
	proxyCmd.Flags().StringVar(&kubeConfigPath, "kubeConfigPath", "./configs/kube/config", "kube config ")

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

	kubernetes.InitClientgo(kubeConfigPath)
}
