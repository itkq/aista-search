package test

import (
	"aista-search/db"
	"encoding/json"
	"github.com/k0kubun/pp"
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

func initEpisodes() {
	_, err := db.Get().Exec("TRUNCATE TABLE episodes")
	if err != nil {
		panic(err)
	}
}

func TestCreate(t *testing.T) {
	initEpisodes()

	// Create episode
	body := httpPost(
		ts.URL+"/api/episode/create",
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
		ts.URL+"/api/episode/create",
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

func TestGetLatest(t *testing.T) {
	initEpisodes()

	// Create episode
	httpPost(
		ts.URL+"/api/episode/create",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)
	httpPost(
		ts.URL+"/api/episode/create",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "2", "title": "fuga"},
	)

	// Get latest episode
	body := httpGet(
		ts.URL+"/api/episode/latest",
		map[string]string{},
	)

	var actual Response2
	json.Unmarshal(body, &actual)
	ep := actual.Episode

	if ep.Id != 2 || ep.Title != "fuga" || ep.Status != 0 {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestUpdate(t *testing.T) {
	initEpisodes()

	// Create episode
	body := httpPost(
		ts.URL+"/api/episode/create",
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
	httpPost(
		ts.URL+"/api/episode/update",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "fuga", "status": "1"},
	)

	body = httpGet(
		ts.URL+"/api/episode/latest",
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
