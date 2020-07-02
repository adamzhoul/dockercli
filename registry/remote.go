package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/adamzhoul/dockercli/common"
	v1 "k8s.io/api/core/v1"
)

type remoteClient struct {
	addr  string
	agent *AgentConfig
}

// type ContainerInfo struct {
// 	Image       string `json:"image"`
// 	ContainerID string `json:"container_id"`
// 	HostIp      string `json:"host_ip"`
// }

const POD_INFO_URL = "http://%s/api/v1/namespaces/%s/pods/%s"
const POD_INFO_LIST = "http://%s/api/v2/namespaces/%s/pods"

func (r *remoteClient) Init(config string, agent *AgentConfig) error {

	if config == "" {
		return errors.New("config empty")
	}

	// todo: check format; ip + addr

	r.addr = config
	r.agent = agent
	fmt.Println("set remote addr", r.addr)
	return nil
}

// get pod container info
// include: containerImage containerID HostIP
func (r remoteClient) FindPodContainerInfo(cluster string, namespace string, podName string, containerName string) (string, string, string, error) {

	url := fmt.Sprintf(POD_INFO_URL, r.addr, namespace, podName)
	res, err := common.HttpGet(url, map[string]string{"cluster": cluster})
	if err != nil {
		return "", "", "", err
	}

	type httpResponse struct {
		Code    int  `json:"code"`
		Success bool `json:"success"`

		Message string `json:"msg"`
		Data    v1.Pod `json:"data"`
	}
	var response httpResponse
	err = json.Unmarshal(res, &response)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}
	if !response.Success {
		log.Println(response.Message)
		return "", "", "", errors.New(response.Message)
	}

	return extraceContainerInfoFromPod(&response.Data, containerName)
}

func (r remoteClient) FindAgentIp(cluster string, hostIP string) (string, error) {

	if r.agent.Ip != "" {
		return r.agent.Ip, nil
	}

	podUrl, _ := url.Parse(fmt.Sprintf(POD_INFO_LIST, r.addr, r.agent.Namespace))
	q, _ := url.ParseQuery(podUrl.RawQuery)
	q.Set("labelSelector", r.agent.Label)
	podUrl.RawQuery = q.Encode()

	res, err := common.HttpGet(podUrl.String(), map[string]string{"cluster": cluster})
	if err != nil {
		return "", err
	}
	type httpResponse struct {
		Code    int  `json:"code"`
		Success bool `json:"success"`

		Message string    `json:"msg"`
		Data    []*v1.Pod `json:"data"`
	}
	var response httpResponse
	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", err
	}
	if !response.Success {
		return "", errors.New(response.Message)
	}

	for _, agent := range response.Data {

		if agent.Status.HostIP == hostIP {
			return agent.Status.PodIP, nil
		}
	}
	return "", errors.New("agent not found")
}

func (r remoteClient) FindAgentPort() string {
	return r.agent.Port
}
