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

type Response2 struct {
	Status  string     `json:"status"`
	Episode db.Episode `json:"episode"`
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
	endpoint := ts.URL + "/api/episode/create"
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

func TestLatest(t *testing.T) {
	_, err := db.Get().Exec("TRUNCATE TABLE episodes")
	if err != nil {
		panic(err)
	}

	router := route.New()
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create episode
	endpoint := ts.URL + "/api/episode/create"
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

	endpoint2 := ts.URL + "/api/episode/latest"
	body = httpGet(
		endpoint2,
		map[string]string{},
	)

	var actual2 Response2
	json.Unmarshal(body, &actual2)
	ep := actual2.Episode

	if ep.Id != 1 || ep.Title != "hoge" || ep.Status != 0 {
		pp.Println(actual2)
		t.Error("response error")
	}
}

func TestUpdate(t *testing.T) {
	_, err := db.Get().Exec("TRUNCATE TABLE episodes")
	if err != nil {
		panic(err)
	}

	router := route.New()
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create episode
	endpoint := ts.URL + "/api/episode/create"
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
	endpoint2 := ts.URL + "/api/episode/update"
	httpPost(
		endpoint2,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "fuga", "status": "1"},
	)

	endpoint3 := ts.URL + "/api/episode/latest"
	body = httpGet(
		endpoint3,
		map[string]string{},
	)

	var actual2 Response2
	json.Unmarshal(body, &actual2)
	ep := actual2.Episode

	if ep.Id != 1 || ep.Title != "fuga" || ep.Status != 1 {
		pp.Println(actual2)
		t.Error("response error")
	}
}
