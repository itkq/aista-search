package test

import (
	"aista-search/db"
	"encoding/json"
	"github.com/k0kubun/pp"
	"testing"
)

type EpisodeResponse struct {
	Status  string     `json:"status"`
	ID      int        `json:"id"`
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

	var actual, expected EpisodeResponse

	// Create episode
	body := httpPost(
		ts.URL+"/api/episode/create",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	json.Unmarshal(body, &actual)
	expected = EpisodeResponse{Status: "ok", ID: db.EpCreated}

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

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	expected = EpisodeResponse{Status: "bad", ID: 0}

	if actual != expected {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestGetEpisode(t *testing.T) {
	initEpisodes()

	var actual, expected EpisodeResponse

	// Get no episode
	body := httpGet(
		ts.URL+"/api/episode/latest",
		map[string]string{},
	)

	json.Unmarshal(body, &actual)
	expected = EpisodeResponse{Status: "bad"}
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

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)

	ep := actual.Episode
	if ep.ID != 2 || ep.Title != "fuga" || ep.Status != db.EpCreated {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get all episodes
	ptr, _ := db.GetEpisodes()
	episodes := *ptr
	ep1 := episodes[0]
	if ep1.ID != 1 || ep1.Title != "hoge" {
		pp.Println(actual)
		t.Error("db error")
	}

	ep2 := episodes[1]
	if ep2.ID != 2 || ep2.Title != "fuga" {
		pp.Println(actual)
		t.Error("db error")
	}
}

func TestUpdateEpisode(t *testing.T) {
	initEpisodes()

	var actual, expected EpisodeResponse

	// Create episode
	body := httpPost(
		ts.URL+"/api/episode/create",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	json.Unmarshal(body, &actual)
	expected = EpisodeResponse{Status: "ok", ID: 1}

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

	json.Unmarshal(body, &actual)
	ep := actual.Episode

	if ep.ID != 1 || ep.Title != "fuga" || ep.Status != db.EpCreated {
		pp.Println(actual)
		t.Error("response error")
	}
}
