package proxy

import (
	"net/http"
)

func handleExec(w http.ResponseWriter, req *http.Request) {

	proxy2Agent(w, req, "/api/v1/exec")

}
