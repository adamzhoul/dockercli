default: build

mac_build:
	GOOS=darwin GOPROXY=https://goproxy.cn go build -o debugctl main.go 

build:
	GOPROXY=https://goproxy.cn go build -o debugctl main.go 

agent:
	./debugctl agent

debug:
	GOPROXY=https://goproxy.cn go build -o debugctl main.go
	./debugctl mock  --addr 0.0.0.0:8083 &
	./debugctl agent --listenAddress 0.0.0.0:8084 &
	./debugctl proxy --addr 0.0.0.0:8082 --sidecar 127.0.0.1:8083 --agp 8084 --registry remote --registryConfig 127.0.0.1:8083

debug_auth:
	./debugctl mock  --addr 0.0.0.0:8083 &
	./debugctl agent --listenAddress 0.0.0.0:8084 &
	./debugctl proxy --addr 0.0.0.0:8082 --sidecar 127.0.0.1:8083 --agp 8084 --registry remote --registryConfig 127.0.0.1:8083

proxy:
	./debugctl proxy --addr 0.0.0.0:18080 --agn ratel --agl "app=webide-agent" --registry remote --registryConfig vk-shark.ccp --sidecar 127.0.0.1:8083

# stopall:
#     ps aux|grep 'debugctl'|grep -v grep|awk '{print "kill -9 "$2}'|sh

forward: 
	sudo kubectl port-forward `kubectl get pod -n ops-system |grep ladder|awk '{print $$1}'` -n ops-system 80:8080



# http://127.0.0.1/exec/cluster/cc/ns/namespace/pod/sky-ladder-prod-748595df85-gvfg4/container/application
# cd /root && ./ladder agent --listenAddress 0.0.0.0:20077
TAG := $(shell date +%Y%m%d%H%M%S)
KIMG := docker.io/xxxx/public/sky-ladder
quick:
	GOPROXY=https://goproxy.cn GOOS=linux go build -o ladder main.go
	docker build -t ${KIMG}:${TAG} -f Dockerfile.kind .
	docker push ${KIMG}:${TAG} 
	cd deploy && kustomize edit set image proxyImg=${KIMG}:${TAG}
	kustomize build deploy | kubectl apply -f -
	# docker cp ladder test-worker2:/root
	kubectl get pod -n ops-system |grep sky|grep -v Running