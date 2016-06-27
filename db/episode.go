package db

import (
	"time"
)

type Episode struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Status    uint      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func GetEpisodes() (*[]Episode, error) {
	var episodes []Episode
	if _, err := dbMap.Select(&episodes, "SELECT * FROM episodes"); err != nil {
		return nil, err
	}

	return &episodes, nil
}

func CreateEpisode(id int, title string, status uint) (*Episode, error) {
	episode := Episode{
		Id:        id,
		Title:     title,
		Status:    status,
		CreatedAt: time.Now(),
	}

	if err := dbMap.Insert(&episode); err != nil {
		return nil, err
	}

	return &episode, nil
}
