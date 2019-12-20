package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func CleanContainer(id string) {

	//log.Println("prepare clean process ->", client.ClientVersion())
	ctx, cancel := context.WithTimeout(context.Background(), runtimeConfig.GracefulExitTimeout) // graceful exit timeout
	defer cancel()

	log.Println("clean container, wait for exit or timeout")
	// ContainerWati will return immediately
	// but, will hang on  errCh channel or statusCodeCh,
	// ctx.timeout will lead errCh get data
	statusCodeCh, errCh := client.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	var rmErr error
	select {
	case err := <-errCh:
		if err != nil {
			log.Println("error waiting container exit, kill with --force.", err)
		}
		rmErr = rmContainer(id, true)
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

	ctx, cancel := context.WithTimeout(context.Background(), runtimeConfig.GracefulExitTimeout) // graceful exit timeout
	defer cancel()
	err := client.ContainerRemove(ctx, id,
		types.ContainerRemoveOptions{
			Force: force,
		})
	if err != nil {
		return err
	}

	return nil
}
