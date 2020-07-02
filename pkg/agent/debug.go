package agent

import (
	"log"
	"net/http"

	"github.com/adamzhoul/dockercli/pkg/docker"

	remoteapi "k8s.io/apimachinery/pkg/util/remotecommand"
	kubeletremote "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
)

// handle attach spdy connection and attach to container
// 1. start sidecar container, share namespaces
// 2. attach to sidecar container
func (s *HTTPAgentServer) handleDebug(w http.ResponseWriter, req *http.Request) {
	if !auth(req){
		http.Error(w, "Unauthorized", 401)
		return
	}

	log.Println("handle debug")
	debugContainerID := req.FormValue("debugContainerID")
	attachImage := req.FormValue("attachImage")
	//debugContainerCmd := req.FormValue("debugContainerCmd")

	var attachTargetContainerID string
	if testAttachTargetContainerID == "" {
		resp, err := docker.CreateContainer(attachImage, debugContainerID)
		if err != nil {
			ResponseErr(w, err, 400)
			return
		}
		log.Println("debug container with image ----->", attachImage, resp.ID)
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
		s.RuntimeConfig.StreamIdleTimeout, // idle timeout will lead server send fin package
		s.RuntimeConfig.StreamCreationTimeout,
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
