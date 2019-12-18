buildandrun: build proxy

mac_build:
	GOOS=darwin GOPROXY=https://goproxy.cn go build -o debugctl main.go 

build:
	GOPROXY=https://goproxy.cn go build -o debugctl main.go 

agent:
	./debugctl agent

proxy:
	./debug proxys