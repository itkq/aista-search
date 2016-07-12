package test

import (
	"aista-search/db"
	"encoding/json"
	"github.com/k0kubun/pp"
	"strconv"
	"testing"
)

type EpisodeResponse struct {
	Status   string       `json:"status"`
	ID       int          `json:"id"`
	Episode  db.Episode   `json:"episode"`
	Episodes []db.Episode `json:"episodes"`
	Message  string       `json:"msg"`
}

func initEpisodes() {
	_, err := db.Get().Exec("TRUNCATE TABLE episodes")
	if err != nil {
		panic(err)
	}
}

func TestCreateEpisode(t *testing.T) {
	initEpisodes()

	var actual EpisodeResponse
	var body []byte

	// Create episode
	body = httpPost(
		ts.URL+"/api/episodes/",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	json.Unmarshal(body, &actual)

	if actual.ID != 1 || actual.Status != "ok" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Check unique episode
	body = httpPost(
		ts.URL+"/api/episodes/",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "fuga"},
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)

	if actual.ID != 0 || actual.Status != "bad" {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestGetEpisode(t *testing.T) {
	initEpisodes()

	var actual EpisodeResponse
	var body []byte

	// Get no episode
	body = httpGet(
		ts.URL+"/api/episodes/",
		map[string]string{},
	)

	json.Unmarshal(body, &actual)
	if len(actual.Episodes) != 0 {
		pp.Println(actual)
		t.Error("response error")
	}

	ids := []int{1, 2}
	titles := []string{"foo", "bar"}

	// Create episode
	for i, _ := range ids {
		httpPost(
			ts.URL+"/api/episodes/",
			map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			map[string]string{"id": strconv.Itoa(ids[i]), "title": titles[i]},
		)
	}

	// Get episodes
	body = httpGet(
		ts.URL+"/api/episodes/",
		map[string]string{},
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	for i, e := range actual.Episodes {
		if e.ID != ids[i] || e.Title != titles[i] || e.Status != db.EpCreated {
			pp.Println(actual)
			t.Error("response error")
		}
	}

	// Get episode
	body = httpGet(
		ts.URL+"/api/episodes/1",
		map[string]string{},
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "ok" || actual.Episode.ID != 1 {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get no episode
	body = httpGet(
		ts.URL+"/api/episodes/3",
		map[string]string{},
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" || actual.Episode.ID != 0 {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestUpdateEpisode(t *testing.T) {
	initEpisodes()

	var actual EpisodeResponse
	var body []byte

	// Create episode
	httpPost(
		ts.URL+"/api/episodes/",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"id": "1", "title": "hoge"},
	)

	// Update episode
	body = httpPut(
		ts.URL+"/api/episodes/1",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"title": "fuga", "status": strconv.Itoa(db.EpRetrieved)},
	)
	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "ok" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get episode
	body = httpGet(
		ts.URL+"/api/episodes/1",
		map[string]string{},
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Episode.Title != "fuga" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Update episode error
	body = httpPut(
		ts.URL+"/api/episodes/2",
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		map[string]string{"title": "fuga", "status": strconv.Itoa(db.EpRetrieved)},
	)
	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get episode
	body = httpGet(
		ts.URL+"/api/episodes/2",
		map[string]string{},
	)
	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" || actual.Episode.ID != 0 {
		pp.Println(actual)
		t.Error("response error")
	}
}
