package docker

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func CleanContainer(id string) {

	log.Println("prepare clean process ->", client.ClientVersion())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // graceful exit timeout
	defer cancel()
	// wait the container gracefully exit
	statusCodeCh, errCh := client.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	log.Println("wait done")
	var rmErr error
	select {
	case err := <-errCh:
		if err != nil {
			log.Println("error waiting container exit, kill with --force", err)
			// timeout or error occurs, try force remove anywawy
			rmErr = rmContainer(id, true)
		}
	case <-statusCodeCh:
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
		log.Println(err)
		return err
	}

	return nil
}
