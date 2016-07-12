package db

import (
	"gopkg.in/guregu/null.v3"
	"log"
	"strings"
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

type Images []Image

const (
	ImagesPerPage = 30
)

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

func GetImageByID(id int) (*Image, error) {
	var image Image
	query := "SELECT * FROM images WHERE id=?"
	if err := dbMap.SelectOne(&image, query, id); err != nil {
		return nil, err
	}

	return &image, nil
}

func GetImages(
	episodeID null.Int,
	sentence null.String,
	toSearch bool,
	toUpload bool,
	cnt null.Int,
) (*[]Image, error) {
	var images []Image
	var query string
	var values []interface{}
	var wheres []string
	var order string

	if sentence.Valid && sentence.String != "" {
		wheres = append(wheres, "sentence like ?")
		values = append(values, "%"+sentence.String+"%")
	}
	if episodeID.Valid && episodeID.Int64 != 0 {
		wheres = append(wheres, "episode_id=?")
		values = append(values, int(episodeID.Int64))
	}
	if toSearch {
		wheres = append(wheres, "sentence IS NOT NULL")
	}
	if toUpload {
		wheres = append(wheres, "url IS NULL")
	}

	order = " ORDER BY episode_id DESC, id DESC "
	if cnt.Valid && cnt.Int64 != 0 {
		order += " LIMIT ?"
		values = append(values, int(cnt.Int64))
	}

	query = "SELECT * FROM images "
	if len(wheres) > 0 {
		query += "WHERE "
		query += strings.Join(wheres, " AND ")
	}
	query += order

	log.Println(query)
	log.Println(values)

	if _, err := dbMap.Select(&images, query, values...); err != nil {
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

func (is *Images) Interface() []interface{} {
	imagesInterface := make([]interface{}, len(*is))
	for i, v := range *is {
		imagesInterface[i] = v
	}

	return imagesInterface
}
