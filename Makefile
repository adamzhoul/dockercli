buildall: server cli_exec

server: 
	GOPROXY=https://goproxy.cn go build -o server  cmd/main.go 

cli_exec: 
	GOPROXY=https://goproxy.cn go build -o exec cmd/exec/exec.go 

run: 
	./server