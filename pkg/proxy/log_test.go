package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestReverseProxy(t *testing.T) {

	// 准备http server
	var testServeMux http.ServeMux
	testServeMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get request:", r.URL)
		w.WriteHeader(http.StatusOK)
	})
	testServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("unexpect request:", r.URL)
		w.WriteHeader(http.StatusOK)
	})
	mockServer := httptest.NewServer(&testServeMux)

	// 准备request
	req := httptest.NewRequest("GET", "http://abcd.whatever.com/abcd/whatever", nil)

	// 准备response
	w := httptest.NewRecorder()

	// 准备反向代理
	uri, _ := url.Parse(mockServer.URL)
	params := url.Values{}
	params.Add("debugContainerID", "containerID")
	uri.RawQuery = params.Encode()
	apiPath := "/healthz"

	// 发起请求
	reverseProxy(w, req, uri, apiPath)
}
