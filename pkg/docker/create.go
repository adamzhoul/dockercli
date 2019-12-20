package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
)

func CreateContainer(image string, targetId string) (*container.ContainerCreateCreatedBody, error) {

	if !strings.HasPrefix(targetId, dockerContainerPrefix) {
		return nil, errors.New("not docker container")
	}

	dockerContainerId := targetId[len(dockerContainerPrefix):]
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(fmt.Sprintf("container:%s", dockerContainerId)),
		UsernsMode:  container.UsernsMode(fmt.Sprintf("container:%s", dockerContainerId)),
		IpcMode:     container.IpcMode(fmt.Sprintf("container:%s", dockerContainerId)),
		PidMode:     container.PidMode(fmt.Sprintf("container:%s", dockerContainerId)),
		CapAdd:      strslice.StrSlice([]string{"SYS_PTRACE", "SYS_ADMIN"}),
	}

	ctx := context.Background()
	resp, err := client.ContainerCreate(ctx, &container.Config{
		Image:      image,
		Entrypoint: []string{"/bin/bash"},
		// Cmd:       []string{"/usr/sbin/adduser --gecos '' --disabled-password coder", "&&", "/bin/bash"},
		//Cmd:       []string{"-c", "/usr/sbin/adduser  --gecos '' --disabled-password coder && /bin/bash"},
		Tty:       true,
		OpenStdin: true,
	}, hostConfig, nil, "")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &resp, nil
}

func RunContainer(id string) error {

	ctx := context.Background()
	err := client.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
