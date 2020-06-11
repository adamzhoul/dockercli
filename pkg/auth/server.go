package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	authApi string
)

func getAuth(token string, resource string) (username string, actions string) {

	if authApi == "" {
		log.Println("authApi empty ,skip ")
		username = "testuser"
		actions = "exec,log,debug"
		return
	}

	// spell request params
	authApiUrl, err := url.Parse(authApi)
	req := &http.Request{
		URL:    authApiUrl,
		Header: http.Header{},
	}
	c := &http.Cookie{
		Name:  "ssoToken",
		Value: token,
	}
	req.AddCookie(c)
	q, _ := url.ParseQuery(authApiUrl.RawQuery)
	q.Set("resource", resource)
	authApiUrl.RawQuery = q.Encode()

	// request auth server
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("respnse fro auth server:", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read auth server resp body", err)
		return
	}

	// parse response data
	type Result struct {
		Code int64 `json:"code"`

		Success bool              `json:"success"` // true 或者 false 代表请求是否成功
		Message string            `json:"message"`
		Data    map[string]string `json:"data"`
	}
	authinfo := Result{}
	err = json.Unmarshal(body, &authinfo)
	if err != nil {
		log.Println("parse auth server resp data ", string(body), err)
		return
	}
	if !authinfo.Success {
		log.Println(authinfo.Message)
		return
	}
	if _, ok := authinfo.Data["username"]; ok {
		username = authinfo.Data["username"]
	}
	if _, ok := authinfo.Data["actions"]; ok {
		actions = authinfo.Data["actions"]
	}

	return
}
