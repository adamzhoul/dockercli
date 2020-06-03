package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adamzhoul/dockercli/common"
	"github.com/adamzhoul/dockercli/registry"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// mock registry client, response specific podinfo etc...
// used for test
var (
	mockAddr string

	containerID string
	image       string
	hostIp      string
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

	mockCmd.Flags().StringVar(&containerID, "cid", "1", "containerID")
	mockCmd.Flags().StringVar(&image, "img", "2", "imageID")
	mockCmd.Flags().StringVar(&hostIp, "ip", "127.0.0.1", "physical machine IP")

}

func route() *mux.Router {

	//mux := http.NewServeMux()
	route := mux.NewRouter()

	route.HandleFunc("/", Index)
	route.HandleFunc("/api/cluster/{cluster}/ns/{namespace}/podname/{podname}/containername/{containername}", podInfo)
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

	i := registry.ContainerInfo{
		Image:       image,
		ContainerID: containerID,
		HostIp:      hostIp,
	}
	d, _ := json.Marshal(i)

	r := common.HttpResponse{
		Code: 0,
		Msg:  "",
		Data: string(d),
	}
	resp, _ := json.Marshal(r)

	w.Write(resp)
}
