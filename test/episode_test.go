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

func TestCreateisode(t *testing.T) {
	initEpisodes()

	var actual EpisodeResponse
	var body []byte

	// Create isode
	body = httpRequest(
		"POST",
		ts.URL+"/api/episodes/",
		&map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		&map[string]string{"id": "1", "title": "hoge"},
	)

	json.Unmarshal(body, &actual)

	if actual.ID != 1 || actual.Status != "ok" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Check unique isode
	body = httpRequest(
		"POST",
		ts.URL+"/api/episodes/",
		&map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		&map[string]string{"id": "1", "title": "fuga"},
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)

	if actual.ID != 0 || actual.Status != "bad" {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestGetisode(t *testing.T) {
	initEpisodes()

	var actual EpisodeResponse
	var body []byte

	// Get no isode
	body = httpRequest(
		"GET",
		ts.URL+"/api/episodes/",
		nil,
		nil,
	)

	json.Unmarshal(body, &actual)
	if len(actual.Episodes) != 0 {
		pp.Println(actual)
		t.Error("response error")
	}

	ids := []int{1, 2}
	titles := []string{"foo", "bar"}

	// Create isode
	for i, _ := range ids {
		httpRequest(
			"POST",
			ts.URL+"/api/episodes/",
			&map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			&map[string]string{"id": strconv.Itoa(ids[i]), "title": titles[i]},
		)
	}

	// Get isodes
	body = httpRequest(
		"GET",
		ts.URL+"/api/episodes/",
		nil,
		nil,
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	for i, e := range actual.Episodes {
		if e.ID != ids[i] || e.Title != titles[i] || e.Status != db.EpCreated {
			pp.Println(actual)
			t.Error("response error")
		}
	}

	// Get isode
	body = httpRequest(
		"GET",
		ts.URL+"/api/episodes/1",
		nil,
		nil,
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "ok" || actual.Episode.ID != 1 {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get no isode
	body = httpRequest(
		"GET",
		ts.URL+"/api/episodes/3",
		nil,
		nil,
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" || actual.Episode.ID != 0 {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestUpdateisode(t *testing.T) {
	initEpisodes()

	var actual EpisodeResponse
	var body []byte

	// Create isode
	httpRequest(
		"POST",
		ts.URL+"/api/episodes/",
		&map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		&map[string]string{"id": "1", "title": "hoge"},
	)

	// Update isode
	body = httpRequest(
		"PUT",
		ts.URL+"/api/episodes/1",
		&map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		&map[string]string{"title": "fuga", "status": strconv.Itoa(db.EpRetrieved)},
	)
	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "ok" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get isode
	body = httpRequest(
		"GET",
		ts.URL+"/api/episodes/1",
		nil,
		nil,
	)

	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Episode.Title != "fuga" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Update isode error
	body = httpRequest(
		"PUT",
		ts.URL+"/api/episodes/2",
		&map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		&map[string]string{"title": "fuga", "status": strconv.Itoa(db.EpRetrieved)},
	)
	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" {
		pp.Println(actual)
		t.Error("response error")
	}

	// Get isode
	body = httpRequest(
		"GET",
		ts.URL+"/api/episodes/2",
		nil,
		nil,
	)
	actual = EpisodeResponse{}
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" || actual.Episode.ID != 0 {
		pp.Println(actual)
		t.Error("response error")
	}
}
