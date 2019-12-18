package kubernetes

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var k8sClient *kubernetes.Clientset

func InitClientgo(kubeConfigPath string) {
	log.Println("load kube config :", kubeConfigPath)
	conf, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	k8sClient, err = kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
}
