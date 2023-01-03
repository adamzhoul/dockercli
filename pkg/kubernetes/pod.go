package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FindPodByName(namespace string, podName string) *v1.Pod {

	if k8sClient == nil {
		log.Println("k8sclient is nil")
		return nil
	}

	pod, err := k8sClient.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return pod
}

func FindPodsByLabel(namespace string, labelSelector string) []v1.Pod {

	pods, err := k8sClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil
	}

	return pods.Items
}
