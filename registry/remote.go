package registry

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adamzhoul/dockercli/common"
)

type remoteClient struct {
	addr string
}

type ContainerInfo struct {
	Image       string `json:"image"`
	ContainerID string `json:"container_id"`
	HostIp      string `json:"host_ip"`
}

const CONTAINER_INFO_URL = "%s/api/cluster/%s/namespace/%s/podname/%s/containername/%s"

func (r remoteClient) Init(config string) error {

	if config == "" {
		return errors.New("config empty")
	}

	// todo: check format; ip + addr

	r.addr = config
	return nil
}

// get pod container info
// include: containerImage containerID HostIP
func (r remoteClient) FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error) {

	url := fmt.Sprintf(CONTAINER_INFO_URL, r.addr, cluster, namespace, podName, containerName)
	res, err := common.HttpGet(url, nil)
	if err != nil {
		return "", "", "", err
	}

	var c ContainerInfo
	err = json.Unmarshal([]byte(res.Data), &c)
	if err != nil {
		return "", "", "", err
	}

	return c.Image, c.ContainerID, c.HostIp, nil
}
