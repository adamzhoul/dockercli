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
