buildall: server client

server: 
	GOPROXY=https://goproxy.cn go build -o ser  cmd/proxy/proxy.go 

client: 
	GOPROXY=https://goproxy.cn go build -o dc cmd/cli/cli.go 

agent:
	GOPROXY=https://goproxy.cn go build -o ag cmd/agent/agent.go 

run: 
	./server