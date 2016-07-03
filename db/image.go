package db

import (
	"database/sql"
	"time"
)

type Image struct {
	Id        int            `db:"id" json:"id"`
	EpisodeId int            `db:"episode_id" json:"episode_id"`
	Path      string         `db:"path" json:"path"`
	Url       sql.NullString `db:"url" json:"url"`
	Sentence  sql.NullString `db:"sentence" json:"sentence"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
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

func GetImagesByEpisodeId(episode_id int) (*[]Image, error) {
	var images []Image
	query := "SELECT * FROM images WHERE episode_id=?"
	if _, err := dbMap.Select(&images, query, episode_id); err != nil {
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
		if _, err := tx.Exec(query, img.Url, img.Sentence, time.Now(), img.Path); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
