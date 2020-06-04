package registry

import (
	"errors"
)

type registryClient interface {
	Init(config string, agent *AgentConfig) error
	FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error)
	FindAgentIp(cluster string, hostIP string) (string, error)
	FindAgentPort() string
}

type AgentConfig struct {
	Namespace string
	Label     string
	Ip        string
	Port      string
}

var Client registryClient

func InitClient(registryType string, config string, agent *AgentConfig) error {

	if registryType == "local" { // use inner k8s client implmentation, .kube/config file is required.
		Client = &k8sClient{}
	} else if registryType == "remote" { // use remote k8s client, remote http address is required
		Client = &remoteClient{}
	} else {
		return errors.New("type not found")
	}

	return Client.Init(config, agent)
}
