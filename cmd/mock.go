package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// mock registry client, response specific podinfo etc...
// used for test
var (
	mockAddr string

	containerName string
	containerID   string
	image         string
	hostIp        string
)

var mockCmd = &cobra.Command{
	Use:           "mock",
	Short:         "mock is a command line tool to mock registry resposne ",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runMock,
}

func init() {

	mockCmd.Flags().StringVar(&mockAddr, "addr", "0.0.0.0:8083", "http listen addr")

	mockCmd.Flags().StringVar(&containerName, "name", "application", "containerName")
	mockCmd.Flags().StringVar(&containerID, "id", "11111", "containerID")
	mockCmd.Flags().StringVar(&image, "img", "222222", "imageID")
	mockCmd.Flags().StringVar(&hostIp, "ip", "127.0.0.1", "physical machine IP")

}

func route() *mux.Router {

	//mux := http.NewServeMux()
	route := mux.NewRouter()

	route.HandleFunc("/", Index)
	route.HandleFunc("/api/v1/namespaces/{namespace}/pods/{podname}", podInfo)
	route.HandleFunc("/healthz", Index)

	return route
}

func runMock(cmd *cobra.Command, args []string) error {

	h := http.Server{
		Addr:    mockAddr,
		Handler: route(),
	}

	log.Println("http server:", mockAddr)
	if err := h.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("I'm OK ~"))
}

func podInfo(w http.ResponseWriter, req *http.Request) {

	pathParams := mux.Vars(req)
	name := pathParams["podname"]
	namespace := pathParams["namespace"]

	p := findPodFromCurrentK8s(name, namespace)
	if p == nil {
		p = mockResult()
	}

	type httpResponse struct {
		Code    int  `json:"code"`
		Success bool `json:"success"`

		Message string `json:"msg"`
		Data    v1.Pod `json:"data"`
	}
	r := httpResponse{
		Code:    0,
		Success: true,
		Message: "",
		Data:    *p,
	}
	resp, _ := json.Marshal(r)

	w.Write(resp)
}

func findPodFromCurrentK8s(name, namespace string) *v1.Pod {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("get k8s config err:", err)
		return nil
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println("new k8s clientset err:", err)
		return nil
	}

	pod, err := cs.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("get pod  err:", err)
		return nil
	}

	if pod.Name == "" {
		fmt.Println("get pod empty")
		return nil
	}

	return pod
}

func mockResult() *v1.Pod {
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  containerName,
					Image: image,
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:        containerName,
					ContainerID: containerID,
				},
			},
			HostIP: hostIp,
		},
	}
}
