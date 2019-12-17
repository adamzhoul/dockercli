# dockercli
docker cli study

copied a lot from projects:
1. https://github.com/aylei/kubectl-debug
2. https://github.com/maoqide/kubeutil


## docker exec -ti xxxx /bin/bash
1. start a http server, watting for connection.
2. support for 2 kinds of client: binary from terminal and web browser.

For binary, it starts a SPDY connection.

For web, it starts websocket connection. use xterm.js of course.


