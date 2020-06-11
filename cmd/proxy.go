package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/adamzhoul/dockercli/common"
	"github.com/adamzhoul/dockercli/pkg/agent"
	"github.com/adamzhoul/dockercli/pkg/auth"
	"github.com/adamzhoul/dockercli/pkg/proxy"
	"github.com/adamzhoul/dockercli/registry"
	"github.com/spf13/cobra"
)

var (
	proxyListenAddress string
	registryConfig     string
	registryType       string // where we can get pod info
	sidecar            string

	authApi string // where we can verify userToken and get privilege

	agentC registry.AgentConfig
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
	proxyCmd.Flags().StringVar(&sidecar, "sidecar", "", "sidecar proxy supported")
	proxyCmd.Flags().StringVar(&authApi, "auth", "", "where we can verify userToken and get privilege")

	proxyCmd.Flags().StringVar(&agentC.Namespace, "agn", agent.AGENT_NAMESPACE, "agent namespace")
	proxyCmd.Flags().StringVar(&agentC.Label, "agl", agent.AGENT_LABEL, "agent label")
	proxyCmd.Flags().StringVar(&agentC.Port, "agp", "18080", "agent port")
	proxyCmd.Flags().StringVar(&agentC.Ip, "agip", "", "agent port")

}

func proxyInit() {

	err := registry.InitClient(registryType, registryConfig, &agentC)
	if err != nil {
		log.Fatal(err)
	}
	common.InitHttpProxy(sidecar)
}

func runProxy(cmd *cobra.Command, args []string) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	auth.InitAuth(authApi)

	// start an HttpServer
	proxyInit()
	proxyConfig := proxy.HTTPConfig{
		ListenAddress: proxyListenAddress,
	}
	proxy := proxy.NewHTTPProxyServer(&proxyConfig)
	proxy.Serve(stop)

	return nil
}
