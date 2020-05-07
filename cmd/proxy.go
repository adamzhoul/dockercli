package cmd

import (
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/registry"
	"github.com/adamzhoul/dockercli/pkg/proxy"
	"github.com/spf13/cobra"
)

var (
	agentAddress       string
	proxyListenAddress string
	registryConfig     string
	skipKube           bool
	registry           string   // where we can get pod info
)

var proxyCmd = &cobra.Command{
	Use:           "proxy",
	Short:         "agent is a command line tool which connect to docker",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runProxy,
}

func init() {

	proxyCmd.Flags().StringVar(&agentAddress, "agent", "", "agent ip port")
	proxyCmd.Flags().StringVar(&proxyListenAddress, "addr", "0.0.0.0:80", "http listener")
	proxyCmd.Flags().StringVar(&registry, "registry", "k8s", "connect to k8s apiserver directly" )
	proxyCmd.Flags().StringVar(&registryConfig, "registryConfig", "./configs/kube/config", "kube config ")
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

	// if !skipKube {
	// 	kubernetes.InitClientgo(kubeConfigPath)
	// }

	// 1. load config from file

	// 2. rewrite params  

	// 3. init registry
	initRegistryClient()
}

func initRegistryClient(){
	
	err := registry.InitClient(registry, registryConfig)
	if err != nil {

	}
}
