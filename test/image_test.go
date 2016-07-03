package test

import (
	"aista-search/db"
	"database/sql"
	"encoding/json"
	"github.com/k0kubun/pp"
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
	body := httpPostJSON(
		ts.URL+"/api/image/create",
		jsonBytes,
	)

	var actual ImageResponse
	json.Unmarshal(body, &actual)
	if actual.Count != 2 {
		pp.Println(actual)
		t.Error("response error")
	}
}

func TestGetImages(t *testing.T) {
	initImages()

	body := httpGet(
		ts.URL+"/api/images",
		map[string]string{},
	)
	var actual ImageResponse
	json.Unmarshal(body, &actual)
	if actual.Status != "bad" {
		pp.Println(actual)
		t.Error("response error")
	}

	body = httpGet(
		ts.URL+"/api/images?episode_id=1",
		map[string]string{},
	)
	var actual2 ImageResponse
	json.Unmarshal(body, &actual2)
	if len(actual2.Images) != 0 {
		pp.Println(actual2)
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
	httpPostJSON(
		ts.URL+"/api/image/create",
		jsonBytes,
	)

	body = httpGet(
		ts.URL+"/api/images?episode_id=1",
		map[string]string{},
	)
	var actual3 ImageResponse
	json.Unmarshal(body, &actual3)
	if len(actual3.Images) != 2 {
		pp.Println(actual3)
		t.Error("response error")
	}
}

func TestUpdateImages(t *testing.T) {
	initImages()

	_, err := db.Get().Exec("TRUNCATE TABLE images")
	if err != nil {
		panic(err)
	}

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
	httpPostJSON(
		ts.URL+"/api/image/create",
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

	var request2 []db.Image
	request2 = append(request2, db.Image{
		Path: path1,
		URL: sql.NullString{
			Valid:  true,
			String: urls[path1],
		},
		Sentence: sql.NullString{
			Valid:  true,
			String: sentences[path1],
		},
	})
	request2 = append(request2, db.Image{
		Path: path2,
		URL: sql.NullString{
			Valid:  true,
			String: urls[path2],
		},
		Sentence: sql.NullString{
			Valid:  true,
			String: sentences[path2],
		},
	})

	jsonBytes, _ = json.Marshal(request2)

	// Update images
	httpPostJSON(
		ts.URL+"/api/image/update",
		jsonBytes,
	)

	body := httpGet(
		ts.URL+"/api/images?episode_id="+strconv.Itoa(episodeID),
		map[string]string{},
	)

	var actual ImageResponse
	json.Unmarshal(body, &actual)
	for _, img := range actual.Images {
		if img.URL.String != urls[img.Path] || img.Sentence.String != sentences[img.Path] {
			pp.Println(img)
			t.Error("response error")
		}
	}

}
