package test

import (
	"aista-search/db"
	"encoding/json"
	"github.com/k0kubun/pp"
	"testing"
)

type EpisodeResponse struct {
	Status  string     `json:"status"`
	Id      int        `json:"id"`
	Episode db.Episode `json:"episode"`
	Message string     `json:"msg"`
}

func initEpisodes() {
	_, err := db.Get().Exec("TRUNCATE TABLE episodes")
	if err != nil {
		panic(err)
	}
}

func TestCreateEpisode(t *testing.T) {
	initEpisodes()

	// Create episode
	body := httpPost(
		ts.URL+"/api/episode/create",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	var actual EpisodeResponse
	json.Unmarshal(body, &actual)
	expected := EpisodeResponse{Status: "ok", Id: 1}

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

	var actual2 EpisodeResponse
	json.Unmarshal(body, &actual2)
	expected = EpisodeResponse{Status: "bad", Id: 0}

	if actual2 != expected {
		pp.Println(actual2)
		t.Error("response error")
	}
}

func TestGetEpisode(t *testing.T) {
	initEpisodes()

	// Get no episode
	body := httpGet(
		ts.URL+"/api/episode/latest",
		map[string]string{},
	)

	var actual EpisodeResponse
	json.Unmarshal(body, &actual)
	expected := EpisodeResponse{Status: "bad"}
	if actual != expected {
		pp.Println(actual)
		t.Error("response error")
	}

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
	body = httpGet(
		ts.URL+"/api/episode/latest",
		map[string]string{},
	)

	var actual2 EpisodeResponse
	json.Unmarshal(body, &actual2)

	ep := actual2.Episode
	if ep.Id != 2 || ep.Title != "fuga" || ep.Status != db.EpCreated {
		pp.Println(actual2)
		t.Error("response error")
	}

	// Get all episodes
	ptr, _ := db.GetEpisodes()
	episodes := *ptr
	ep1 := episodes[0]
	if ep1.Id != 1 || ep1.Title != "hoge" {
		pp.Println(actual)
		t.Error("db error")
	}

	ep2 := episodes[1]
	if ep2.Id != 2 || ep2.Title != "fuga" {
		pp.Println(actual)
		t.Error("db error")
	}
}

func TestUpdateEpisode(t *testing.T) {
	initEpisodes()

	// Create episode
	body := httpPost(
		ts.URL+"/api/episode/create",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	var actual EpisodeResponse
	json.Unmarshal(body, &actual)
	expected := EpisodeResponse{Status: "ok", Id: 1}

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

	var actual2 EpisodeResponse
	json.Unmarshal(body, &actual2)
	ep := actual2.Episode

	if ep.Id != 1 || ep.Title != "fuga" || ep.Status != db.EpCreated {
		pp.Println(actual2)
		t.Error("response error")
	}
}
