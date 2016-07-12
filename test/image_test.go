package test

import (
	"aista-search/db"
	"encoding/json"
	"github.com/k0kubun/pp"
	"gopkg.in/guregu/null.v3"
	"strconv"
	"testing"
)

type ImageResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"msg"`
	Count   int        `json:"count"`
	Images  []db.Image `json:"images"`
}

func initImages() {
	_, err := db.Get().Exec("TRUNCATE TABLE images")
	if err != nil {
		panic(err)
	}
}

func TestCreateImages(t *testing.T) {
	initImages()

	var actual ImageResponse
	var body []byte

	var request []db.Image
	request = append(request, db.Image{
		EpisodeID: 1,
		Path:      "img/001/001.jpg",
	})
	request = append(request, db.Image{
		EpisodeID: 1,
		Path:      "img/001/002.jpg",
	})

	jsonBytes, _ := json.Marshal(request)

	// Create images
	body = httpRequestJSON(
		"POST",
		ts.URL+"/api/images/",
		jsonBytes,
	)

	json.Unmarshal(body, &actual)
	if actual.Count != 2 {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestGetImages(t *testing.T) {
	initImages()

	var actual ImageResponse
	var body []byte

	body = httpRequest(
		"GET",
		ts.URL+"/api/images/",
		nil,
		nil,
	)

	json.Unmarshal(body, &actual)
	if actual.Status != "ok" || len(actual.Images) != 0 {
		pp.Println(actual)
		t.Error("response error")
	}

	body = httpRequest(
		"GET",
		ts.URL+"/api/images/?episode_id=1",
		nil,
		nil,
	)

	actual = ImageResponse{}
	json.Unmarshal(body, &actual)
	if len(actual.Images) != 0 {
		pp.Println(actual)
		t.Error("response error")
	}

	var request []db.Image
	request = append(request, db.Image{
		EpisodeID: 1,
		Path:      "img/001/003.jpg",
	})
	request = append(request, db.Image{
		EpisodeID: 1,
		Path:      "img/001/004.jpg",
	})

	jsonBytes, _ := json.Marshal(request)

	// Create images
	httpRequestJSON(
		"POST",
		ts.URL+"/api/images/",
		jsonBytes,
	)

	body = httpRequest(
		"GET",
		ts.URL+"/api/images/?episode_id=1",
		nil,
		nil,
	)

	actual = ImageResponse{}
	json.Unmarshal(body, &actual)
	if len(actual.Images) != 2 {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestUpdateImages(t *testing.T) {
	initImages()

	var actual ImageResponse
	var body []byte

	episodeID := 1
	path1 := "img/001/001.jpg"
	path2 := "img/001/002.jpg"

	var request []db.Image
	request = append(request, db.Image{
		EpisodeID: episodeID,
		Path:      path1,
	})
	request = append(request, db.Image{
		EpisodeID: episodeID,
		Path:      path2,
	})

	jsonBytes, _ := json.Marshal(request)

	// Create images
	httpRequestJSON(
		"POST",
		ts.URL+"/api/images/",
		jsonBytes,
	)

	urls := map[string]string{
		path1: "https://hoge.com/001/001.jpg",
		path2: "https://hoge.com/001/002.jpg",
	}
	sentences := map[string]string{
		path1: "foo",
		path2: "bar",
	}

	request = []db.Image{}
	request = append(request, db.Image{
		Path:     path1,
		URL:      null.StringFrom(urls[path1]),
		Sentence: null.StringFrom(sentences[path1]),
	})
	request = append(request, db.Image{
		Path:     path2,
		URL:      null.StringFrom(urls[path2]),
		Sentence: null.StringFrom(sentences[path2]),
	})

	jsonBytes, _ = json.Marshal(request)

	// Update images
	httpRequestJSON(
		"PUT",
		ts.URL+"/api/images/",
		jsonBytes,
	)

	body = httpRequest(
		"GET",
		ts.URL+"/api/images/?episode_id="+strconv.Itoa(episodeID),
		nil,
		nil,
	)

	json.Unmarshal(body, &actual)
	for _, img := range actual.Images {
		if img.URL.String != urls[img.Path] || img.Sentence.String != sentences[img.Path] {
			pp.Println(img)
			t.Error("response error")
		}
	}

}
