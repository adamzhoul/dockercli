package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
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
	mockCmd.Flags().StringVar(&containerID, "id", "1", "containerID")
	mockCmd.Flags().StringVar(&image, "img", "2", "imageID")
	mockCmd.Flags().StringVar(&hostIp, "ip", "127.0.0.1", "physical machine IP")

}

func route() *mux.Router {

	//mux := http.NewServeMux()
	route := mux.NewRouter()

	route.HandleFunc("/", Index)
	route.HandleFunc("/api/v1/cluster/{cluster}/ns/{namespace}/podname/{podname}", podInfo)
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
	w.Write([]byte("I'm OK!"))
}

func podInfo(w http.ResponseWriter, req *http.Request) {

	p := v1.Pod{
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
		},
	}
	//d, _ := json.Marshal(p)
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
		Data:    p,
	}
	resp, _ := json.Marshal(r)

	w.Write(resp)
}
