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
// 1. start sidecar container, share namespaces
// 2. attach to sidecar container
func handleDebug(w http.ResponseWriter, req *http.Request) {

	log.Println("handle debug")
	debugContainerID := req.FormValue("containerID")
	debugImage := req.FormValue("debugImage")
	//debugContainerCmd := req.FormValue("debugContainerCmd")

	// 1. start sidecar container, with specific image.
	// If debugTargetContainer is empty
	var attachTargetContainerID string
	if testAttachTargetContainerID == "" {
		resp, err := docker.CreateContainer(debugImage, debugContainerID)
		if err != nil {
			log.Println(err)
			return
		}
		err = docker.RunContainer(resp.ID)
		if err != nil {
			log.Println(err)
			return
		}
		attachTargetContainerID = resp.ID
	} else {
		attachTargetContainerID = testAttachTargetContainerID
	}
	defer docker.CleanContainer(attachTargetContainerID)

	// 2. attach to container
	streamOpts := &kubeletremote.Options{
		Stdin:  true,
		Stdout: true,
		Stderr: false,
		TTY:    true,
	}
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
