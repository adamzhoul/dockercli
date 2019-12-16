package kubernetes

// find pod info , node ip included
func FindPod(podName string) {

	if k8sClient == nil {
		return
	}

	// k8sClient.AppsV1().RESTClient().Get()

}
