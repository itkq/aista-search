package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func HttpRequest(
	method string,
	endpoint string,
	header *map[string]string,
	params *map[string]string,
) []byte {
	var req *http.Request
	var values url.Values = url.Values{}
	switch method {
	case "GET":
		req, _ = http.NewRequest(method, endpoint, nil)
	case "PUT", "POST":
		for k, v := range *params {
			values.Add(k, v)
		}
		req, _ = http.NewRequest(method, endpoint, strings.NewReader(values.Encode()))
	default:
		panic("method error")
	}

	if header != nil {
		for k, v := range *header {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return body
}

func HttpRequestJSON(method string, endpoint string, json []byte) []byte {
	if method != "POST" && method != "PUT" {
		panic("method error")
	}
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return body
}
