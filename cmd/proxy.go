package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/pkg/proxy"
	"github.com/adamzhoul/dockercli/registry"
	"github.com/spf13/cobra"
)

var (
	proxyListenAddress string
	registryConfig     string
	registryType       string // where we can get pod info
)

var proxyCmd = &cobra.Command{
	Use:           "proxy",
	Short:         "proxy is a command line tool which connect to agent",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runProxy,
}

func init() {

	proxyCmd.Flags().StringVar(&proxyListenAddress, "addr", "0.0.0.0:80", "http listener")
	proxyCmd.Flags().StringVar(&registryType, "registry", "local", "connect to k8s apiserver directly")
	proxyCmd.Flags().StringVar(&registryConfig, "registryConfig", "./configs/kube/config", "kube config ")

}

func proxyInit() {
	initRegistryClient()
}

func runProxy(cmd *cobra.Command, args []string) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	proxyInit()
	proxyConfig := proxy.HTTPConfig{
		ListenAddress: proxyListenAddress,
	}

	// start an HttpServer
	proxy := proxy.NewHTTPProxyServer(&proxyConfig)
	proxy.Serve(stop)

	return nil
}

func initRegistryClient() {

	err := registry.InitClient(registryType, registryConfig)
	if err != nil {
		log.Fatal(err)
	}
}
