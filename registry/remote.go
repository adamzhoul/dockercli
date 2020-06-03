package registry

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adamzhoul/dockercli/common"
	v1 "k8s.io/api/core/v1"
)

type remoteClient struct {
	addr string
}

// type ContainerInfo struct {
// 	Image       string `json:"image"`
// 	ContainerID string `json:"container_id"`
// 	HostIp      string `json:"host_ip"`
// }

const POD_INFO_URL = "http://%s/api/v1/cluster/%s/namespace/%s/podname/%s"

func (r *remoteClient) Init(config string) error {

	if config == "" {
		return errors.New("config empty")
	}

	// todo: check format; ip + addr

	r.addr = config
	fmt.Println("set remote addr", r.addr)
	return nil
}

// get pod container info
// include: containerImage containerID HostIP
func (r remoteClient) FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error) {

	url := fmt.Sprintf(POD_INFO_URL, r.addr, cluster, namespace, podName)
	res, err := common.HttpGet(url, nil)
	if err != nil {
		return "", "", "", err
	}

	var p v1.Pod
	err = json.Unmarshal([]byte(res.Data), &p)
	if err != nil {
		return "", "", "", err
	}

	return extraceContainerInfoFromPod(&p, containerName)
}
