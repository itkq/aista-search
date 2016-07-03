package db

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type Image struct {
	ID        int         `db:"id" json:"id"`
	EpisodeID int         `db:"episode_id" json:"episode_id"`
	Path      string      `db:"path" json:"path"`
	URL       null.String `db:"url" json:"url"`
	Sentence  null.String `db:"sentence" json:"sentence"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

func CreateImages(images []Image) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	var query string
	for _, img := range images {
		query = "INSERT INTO images (episode_id, path) VALUES (?, ?)"
		if _, err := tx.Exec(query, img.EpisodeID, img.Path); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetImagesByEpisodeID(episodeID int) (*[]Image, error) {
	var images []Image
	query := "SELECT * FROM images WHERE episode_id=?"
	if _, err := dbMap.Select(&images, query, episodeID); err != nil {
		return nil, err
	}

	return &images, nil
}

func UpdateImages(images []Image) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE images SET url=?, sentence=?, updated_at=? WHERE path=?"
	for _, img := range images {
		if _, err := tx.Exec(query, img.URL, img.Sentence, time.Now(), img.Path); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
