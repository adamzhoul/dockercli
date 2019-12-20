package proxy

import "net/http"

func handleLog(w http.ResponseWriter, req *http.Request) {
	proxy2Agent(w, req, "/api/v1/log")
}
