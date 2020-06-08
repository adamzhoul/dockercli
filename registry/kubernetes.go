package registry

import (
	"errors"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type k8sClient struct {
	client *kubernetes.Clientset
	agent  *AgentConfig
	registryClient
}

func (kc *k8sClient) Init(config string, agent *AgentConfig) error {
	log.Println("load kube config :", config)
	conf, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return err
	}

	kc.client, err = kubernetes.NewForConfig(conf)
	if err != nil {
		return err
	}
	kc.agent = agent

	return nil
}

func (kc k8sClient) findPodByName(namespace string, podName string) *v1.Pod {

	if kc.client == nil {
		log.Println("k8sclient is nil")
		return nil
	}

	pod, err := kc.client.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return pod
}

func (kc k8sClient) findPodsByLabel(namespace string, labelSelector string) []v1.Pod {

	pods, err := kc.client.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil
	}

	return pods.Items
}

// get pod container info
// include: containerImage containerID HostIP
func (kc k8sClient) FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error) {

	// 1. find pod
	pod := kc.findPodByName(namespace, podName)
	if pod == nil {
		return "", "", "", errors.New("pod not found")
	}

	return extraceContainerInfoFromPod(pod, containerName)
}

func (kc k8sClient) FindAgentIp(cluster string, hostIP string) (string, error) {
	agents := kc.findPodsByLabel(kc.agent.Namespace, kc.agent.Label)
	for _, agent := range agents {

		if agent.Status.HostIP == hostIP {
			return agent.Status.PodIP, nil
		}
	}

	return "", errors.New("agent not found")
}

func (kc k8sClient) FindAgentPort() string {
	return kc.agent.Port
}

func extraceContainerInfoFromPod(pod *v1.Pod, containerName string) (string, string, string, error) {

	var image, containerID string

	// 1. find container image
	for _, container := range pod.Spec.Containers {

		if container.Name == containerName {
			image = container.Image
			break
		}
	}

	// 2. find container ID
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == containerName {
			containerID = containerStatus.ContainerID
			break
		}
	}

	if image == "" || containerID == "" {
		return image, containerID, pod.Status.HostIP, errors.New("pod info error ")
	}

	return image, containerID, pod.Status.HostIP, nil
}
