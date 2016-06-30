package test

import (
	"aista-search/db"
	"encoding/json"
	"github.com/k0kubun/pp"
	"testing"
)

type ImageResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
	Count   int    `json:"count"`
}

func TestCreateImages(t *testing.T) {
	_, err := db.Get().Exec("TRUNCATE TABLE images")
	if err != nil {
		panic(err)
	}

	var request []db.Image
	request = append(request, db.Image{
		EpisodeId: 1,
		Path:      "img/001/001.jpg",
	})
	request = append(request, db.Image{
		EpisodeId: 1,
		Path:      "img/001/002.jpg",
	})
	request = append(request, db.Image{
		EpisodeId: 1,
		Path:      "img/001/003.jpg",
	})

	jsonBytes, _ := json.Marshal(request)

	// Create images
	body := httpPostJSON(
		ts.URL+"/api/image/create",
		jsonBytes,
	)

	var actual ImageResponse
	json.Unmarshal(body, &actual)
	if actual.Count != 3 {
		pp.Println(actual)
		t.Error("response error")
	}
}
