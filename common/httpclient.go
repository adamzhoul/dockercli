package common

import (
	"io/ioutil"
	"log"
	"net/http"
)

func HttpGet(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	log.Println("HTTP GET", resp.StatusCode, url)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func HttpPost(url string, data []byte, header map[string]string) ([]byte, error) {

	return []byte{}, nil
}
