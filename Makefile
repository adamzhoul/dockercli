buildall: server client

server: 
	GOPROXY=https://goproxy.cn go build -o ser  cmd/server/server.go 

client: 
	GOPROXY=https://goproxy.cn go build -o dc cmd/cli/cli.go 

run: 
	./server