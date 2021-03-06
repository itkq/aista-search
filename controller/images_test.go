package controller

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"gopkg.in/guregu/null.v3"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type ImageResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"msg"`
	Count   int        `json:"count"`
	Images  []db.Image `json:"images"`
}

var imageTs *httptest.Server
var imageToken string

func init() {
	os.Setenv("GO_ENV", "test")
	config.LoadEnv()
	db.Connect()

	_, err := db.Get().Exec("TRUNCATE TABLE tokens")
	if err != nil {
		panic(err)
	}
	imageToken = "098f6bcd4621d373cade4e832627b4f6"
	db.CreateToken(imageToken)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/images/", API.ImagesGET)
		api.POST("/images/", API.ImagesPOST)
		api.PUT("/images/", API.ImagesPUT)
	}
	imageTs = httptest.NewServer(router)
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
	body = test.HttpRequestJSON(
		"POST",
		imageTs.URL+"/api/images/?token="+imageToken,
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

	body = test.HttpRequest(
		"GET",
		imageTs.URL+"/api/images/?token="+imageToken,
		nil,
		nil,
	)

	json.Unmarshal(body, &actual)
	if actual.Status != "ok" || len(actual.Images) != 0 {
		pp.Println(actual)
		t.Error("response error")
	}

	body = test.HttpRequest(
		"GET",
		imageTs.URL+"/api/images/?episode_id=1&token="+imageToken,
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
	test.HttpRequestJSON(
		"POST",
		imageTs.URL+"/api/images/",
		jsonBytes,
	)

	body = test.HttpRequest(
		"GET",
		imageTs.URL+"/api/images/?episode_id=1&token="+imageToken,
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
	test.HttpRequestJSON(
		"POST",
		imageTs.URL+"/api/images/?token="+imageToken,
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
	test.HttpRequestJSON(
		"PUT",
		imageTs.URL+"/api/images/?token="+imageToken,
		jsonBytes,
	)

	body = test.HttpRequest(
		"GET",
		imageTs.URL+"/api/images/?episode_id="+strconv.Itoa(episodeID)+"&token="+imageToken,
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
