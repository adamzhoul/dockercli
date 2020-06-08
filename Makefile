default: build

mac_build:
	GOOS=darwin GOPROXY=https://goproxy.cn go build -o debugctl main.go 

build:
	GOPROXY=https://goproxy.cn go build -o debugctl main.go 

agent:
	./debugctl agent
	
proxy:
	./debugctl proxy --addr 0.0.0.0:18080 --agn ratel --agl "app=webide-agent" --registry remote --registryConfig vk-shark.ccp --sidecar 127.0.0.1:8083