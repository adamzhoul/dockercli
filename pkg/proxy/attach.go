package proxy

import "net/http"

// handle websocket connection
// find container node
// connect to agent which on the same node

func handleDebug(w http.ResponseWriter, req *http.Request) {
	proxy2Agent(w, req, "/api/v1/debug")
}