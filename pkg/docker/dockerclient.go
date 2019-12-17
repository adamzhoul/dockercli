package docker

import (
	dockerclient "github.com/docker/docker/client"
	"log"
)

var client *dockerclient.Client

// init connection to docker.sock
func InitDockerclientConn() {
	c, err := dockerclient.NewClient("unix:///var/run/docker.sock", "", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	client = c
}

func CloseConn() {
	client.Close()
}
