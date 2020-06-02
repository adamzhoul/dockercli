package registry

type registryClient interface{
	Init(config string)
	FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error) 
}

var Client registryClient

func InitClient(registryType string, config string) error{

	if registryType == "k8s"{
		Client = &k8sClient{}
	}

	Client.Init(config)
	return nil
}
