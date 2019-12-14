buildall: server client

server: 
	GOPROXY=https://goproxy.cn go build -o server  cmd/server/server.go 

client: 
	GOPROXY=https://goproxy.cn go build -o cli cmd/cli/cli.go 

run: 
	./server