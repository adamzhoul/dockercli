module github.com/adamzhoul/dockercli

go 1.13

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/elazarl/goproxy v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/segmentio/ksuid v1.0.3
	github.com/spf13/cobra v1.6.0
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/net v0.3.1-0.20221206200815-1e63c2f08a10 // indirect
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	k8s.io/api v0.26.0
	k8s.io/apimachinery v0.26.0
	k8s.io/apiserver v0.22.0
	k8s.io/client-go v0.26.0
	k8s.io/component-base v0.26.0 // indirect
	k8s.io/cri-api v0.23.1
	k8s.io/klog/v2 v2.80.1
	k8s.io/kubernetes v1.22.0
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect

)

//replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1

replace k8s.io/api => k8s.io/api v0.22.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.23.0-alpha.0

replace k8s.io/apiserver => k8s.io/apiserver v0.22.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.0

replace k8s.io/client-go => k8s.io/client-go v0.22.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.0

replace k8s.io/code-generator => k8s.io/code-generator v0.22.2-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.22.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.22.0

replace k8s.io/controller-manager => k8s.io/controller-manager v0.22.0

replace k8s.io/cri-api => k8s.io/cri-api v0.23.0-alpha.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.0

replace k8s.io/kubectl => k8s.io/kubectl v0.22.0

replace k8s.io/kubelet => k8s.io/kubelet v0.22.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.0

replace k8s.io/metrics => k8s.io/metrics v0.22.0

replace k8s.io/mount-utils => k8s.io/mount-utils v0.22.1-rc.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.22.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.0

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.22.0

replace k8s.io/sample-controller => k8s.io/sample-controller v0.22.0
