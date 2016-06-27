package test

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"encoding/json"
	"github.com/k0kubun/pp"
	"net/http/httptest"
	"os"
	"testing"
)

type Response struct {
	Status string `json:"status"`
	Id     int    `json:"id"`
}

func init() {
	os.Setenv("GO_ENV", "test")
	config.LoadEnv()
	db.Connect()
}

func TestCreate(t *testing.T) {
	_, err := db.Get().Exec("TRUNCATE TABLE episodes")
	if err != nil {
		panic(err)
	}

	router := route.New()
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create episode
	endpoint := ts.URL + "/episodes"
	body := httpPost(
		endpoint,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	var actual Response
	json.Unmarshal(body, &actual)
	expected := Response{Status: "ok", Id: 1}

	if actual != expected {
		pp.Println(actual)
		t.Error("response error")
	}

	// Check unique episode
	body = httpPost(
		endpoint,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "fuga"},
	)

	var actual2 Response
	json.Unmarshal(body, &actual2)
	expected = Response{Status: "bad", Id: 0}

	if actual2 != expected {
		pp.Println(actual2)
		t.Error("response error")
	}
}
