package docker

import (
	"fmt"
	dockerclient "github.com/docker/docker/client"
	"log"
)

var client *dockerclient.Client

// init connection to docker.sock
func InitDockerclientConn(dockerAddress string) {
	c, err := dockerclient.NewClient(fmt.Sprintf("unix://%s", dockerAddress), "", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	client = c
}

func CloseConn() {
	client.Close()
}
