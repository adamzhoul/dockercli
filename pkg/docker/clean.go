package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func CleanContainer(id string) {

	ctx := context.Background()
	// wait the container gracefully exit
	statusCode, err := client.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	var rmErr error
	if err != nil {
		log.Println("error waiting container exit, kill with --force")
		// timeout or error occurs, try force remove anywawy
		rmErr = rmContainer(id, true)
	} else {
		log.Println("container return response code:", statusCode)
		rmErr = rmContainer(id, false)
	}

	if rmErr != nil {
		log.Printf("error remove container: %s \n", id)
	} else {
		log.Printf("Debug session end, debug container %s removed", id)
	}
}

func rmContainer(id string, force bool) error {

	ctx := context.Background()
	err := client.ContainerRemove(ctx, id,
		types.ContainerRemoveOptions{
			Force: force,
		})
	if err != nil {
		return err
	}

	return nil
}
