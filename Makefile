buildall: server client

server: 
	GOPROXY=https://goproxy.cn go build -o ser  cmd/proxy/proxy.go 

client: 
	GOPROXY=https://goproxy.cn go build -o dc cmd/cli/cli.go 

run: 
	./server