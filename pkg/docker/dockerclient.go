package docker

import (
	"fmt"
	"log"
	"time"

	dockerclient "github.com/docker/docker/client"
)

const (
	dockerContainerPrefix = "docker://"
)

type RuntimeConfig struct {
	DockerEndpoint        string
	GracefulExitTimeout   time.Duration
	StreamIdleTimeout     time.Duration
	StreamCreationTimeout time.Duration
}

var runtimeConfig RuntimeConfig
var client *dockerclient.Client

// init connection to docker.sock
func InitDockerclientConn(config RuntimeConfig) {
	c, err := dockerclient.NewClient(fmt.Sprintf("unix://%s", config.DockerEndpoint), "", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	client = c
	runtimeConfig = config
}

func CloseConn() {
	client.Close()
}
