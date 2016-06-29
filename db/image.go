package db

import (
	"time"
)

type Image struct {
	Id        int       `db:"id" json:"id"`
	EpisodeId int       `db:"episode_id" json:"episode_id"`
	Path      string    `db:"path" json:"path"`
	Url       string    `db:"url" json:"url"`
	Sentence  string    `db:"sentence" json:"sentence"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func NewImage(episodeId int, path string) Image {
	return Image{EpisodeId: episodeId, Path: path}
}

func CreateImages(images []Image) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	var query string
	for _, img := range images {
		query = "INSERT INTO images (episode_id, path) VALUES (?, ?)"
		if _, err := tx.Exec(query, img.EpisodeId, img.Path); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
