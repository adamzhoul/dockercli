package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// // find pods info , node ip included
// func FindPodsByName(namespace string, podsName []string) []v1.Pod {

// 	if k8sClient == nil {
// 		return nil
// 	}

// 	pods, err := k8sClient.CoreV1().Pods("namespace").List(metav1.ListOptions{})
// 	if err != nil {
// 		// handle error
// 		return nil
// 	}

// 	res := []v1.Pod{}
// 	for _, pod := range pods.Items {
// 		fmt.Println(pod.Name, pod.Status.PodIP, pod.Spec.NodeName)
// 		res = append(res, pod)
// 	}

// 	return res

// }

func FindPodByName(namespace string, podName string) *v1.Pod {

	if k8sClient == nil {
		return nil
	}

	pod, err := k8sClient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil
	}

	return pod
}

func FindPodsByLabel(namespace string, labelSelector string) []v1.Pod {

	pods, err := k8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil
	}

	return pods.Items
}
