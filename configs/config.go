package configs

type AgentConfig struct {
	Ip   string
	Port int
}

type ProxyConfig struct {
	Port  int
	Agent AgentConfig
}

var defaultAgentConfig = AgentConfig{
	Ip:   "",
	Port: 80,
}

var defaultProxyConfig = ProxyConfig{
	Port:  80,
	Agent: defaultAgentConfig,
}

var defaultKubeConfigPath = "./config/kube/kube-config"

func InitConfigs() {

}
