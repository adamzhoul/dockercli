package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func HttpGet(url string, header map[string]string) (HttpResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	for k, v := range header {
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return HttpResponse{}, err
	}
	defer resp.Body.Close()
	log.Println("GET", resp.StatusCode, url)

	var r HttpResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return HttpResponse{}, err
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return HttpResponse{}, err
	}

	return r, nil
}

func HttpPost(url string, data []byte, header map[string]string) (HttpResponse, error) {

	return HttpResponse{}, nil
}
