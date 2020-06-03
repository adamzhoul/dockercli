package registry

import "errors"

type registryClient interface {
	Init(config string) error
	FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error)
}

var Client registryClient

func InitClient(registryType string, config string) error {

	if registryType == "local" { // use inner k8s client implmentation, .kube/config file is required.
		Client = &k8sClient{}
	} else if registryType == "remote" { // use remote k8s client, remote http address is required
		Client = &remoteClient{}
	} else {
		return errors.New("type not found")
	}

	return Client.Init(config)
}
