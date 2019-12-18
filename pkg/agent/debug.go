package agent

import (
	"log"
	"net/http"
	"time"

	"github.com/adamzhoul/dockercli/pkg/docker"

	remoteapi "k8s.io/apimachinery/pkg/util/remotecommand"
	kubeletremote "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
)

// handle attach spdy connection and attach to container
// 1. upgrade connection to spdy
// 2. find target container
// 3. pull image
// 4. start sidecar container, share namespaces
// 5. attach to sidecar container
func handleDebug(w http.ResponseWriter, req *http.Request) {

	log.Println("handle debug")
	streamOpts := &kubeletremote.Options{
		Stdin:  true,
		Stdout: true,
		Stderr: false,
		TTY:    true,
	}
	attachTargetContainerID := attachDebugTargetContainerID
	kubeletremote.ServeAttach(
		w,
		req,
		GetAttacher(),
		"",
		"",
		attachTargetContainerID,
		streamOpts,
		10*time.Minute,
		15*time.Second,
		remoteapi.SupportedStreamingProtocols)

}

// get attacher ,who do the attach work
func GetAttacher() *docker.ContainerAttacher {

	return docker.NewContainerAttacher()
}

// func pullImage(image string, client *dockerclient.Client) {

// 	//authBytes := base64.URLEncoding.EncodeToString([]byte(authStr))
// 	// types.AuthConfig{}
// 	ctx := context.Background()
// 	out, err := client.ImagePull(ctx, image, types.ImagePullOptions{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(out)
// 	defer out.Close()

// 	io.Copy(os.Stdout, out)
// }

// func createContainer(image string, targetId string, client *dockerclient.Client) (*container.ContainerCreateCreatedBody, error) {

// 	hostConfig := &container.HostConfig{
// 		NetworkMode: container.NetworkMode(fmt.Sprintf("container:%s", targetId)),
// 		UsernsMode:  container.UsernsMode(fmt.Sprintf("container:%s", targetId)),
// 		IpcMode:     container.IpcMode(fmt.Sprintf("container:%s", targetId)),
// 		PidMode:     container.PidMode(fmt.Sprintf("container:%s", targetId)),
// 		CapAdd:      strslice.StrSlice([]string{"SYS_PTRACE", "SYS_ADMIN"}),
// 	}

// 	ctx := context.Background()
// 	resp, err := client.ContainerCreate(ctx, &container.Config{
// 		Image:      image,
// 		Entrypoint: []string{"/bin/bash"},
// 		// Cmd:       []string{"/usr/sbin/adduser --gecos '' --disabled-password coder", "&&", "/bin/bash"},
// 		Cmd:       []string{"-c", "/usr/sbin/adduser  --gecos '' --disabled-password coder && /bin/bash"},
// 		Tty:       true,
// 		OpenStdin: true,
// 	}, hostConfig, nil, "")
// 	if err != nil {
// 		panic(err)
// 	}

// 	return &resp, nil
// }

// func runContainer(id string, client *dockerclient.Client) error {

// 	ctx := context.Background()
// 	err := client.ContainerStart(ctx, id, types.ContainerStartOptions{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil
// }

// func cleanContainer(id string, client *dockerclient.Client) {

// 	ctx := context.Background()
// 	// wait the container gracefully exit
// 	statusCode, err := client.ContainerWait(ctx, id)
// 	var rmErr error
// 	if err != nil {
// 		log.Println("error waiting container exit, kill with --force")
// 		// timeout or error occurs, try force remove anywawy
// 		rmErr = rmContainer(id, true, client)
// 	} else {
// 		log.Println("container return response code:", statusCode)
// 		rmErr = rmContainer(id, false, client)
// 	}

// 	if rmErr != nil {
// 		log.Printf("error remove container: %s \n", id)
// 	} else {
// 		log.Printf("Debug session end, debug container %s removed", id)
// 	}

// }

// func rmContainer(id string, force bool, client *dockerclient.Client) error {

// 	ctx := context.Background()
// 	err := client.ContainerRemove(ctx, id,
// 		types.ContainerRemoveOptions{
// 			Force: force,
// 		})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
