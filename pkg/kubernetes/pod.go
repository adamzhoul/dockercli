package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// find pods info , node ip included
func FindPod(podsName []string) {

	if k8sClient == nil {
		return
	}

	// k8sClient.AppsV1().RESTClient().Get()
	pods, err := k8sClient.CoreV1().Pods("namespace").List(metav1.ListOptions{})
	if err != nil {
		// handle error
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Status.PodIP)
	}

}
