package test

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
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

func httpGet(endpoint string, header map[string]string) []byte {
	req, _ := http.NewRequest("GET", endpoint, nil)
	for k, v := range header {
		req.Header.Set(k, v)
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

func httpPost(endpoint string, header map[string]string, params map[string]string) []byte {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
	for k, v := range header {
		req.Header.Set(k, v)
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
