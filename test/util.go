package test

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
)

var ts *httptest.Server

func init() {
	os.Setenv("GO_ENV", "test")
	config.LoadEnv()
	db.Connect()

	router := route.New()
	ts = httptest.NewServer(router)
}

func httpRequest(
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

func httpPostJSON(endpoint string, json []byte) []byte {
	req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(json))
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
