package proxy

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestExtraceResource(t *testing.T) {

	req := httptest.NewRequest("GET", "http://127.0.0.1/api/v1/log/cluster/cc/ns/namespace/ops-system/pod/abc-xxx-def-75575f479f-l4wfw/container/xxx-def", nil)

	r, a := extractResourceActionFromUrl(req)
	fmt.Println("get resource:", r, a)
}
