package db

import (
	"gopkg.in/guregu/null.v3"
	"log"
	"os"
	"os/exec"
	"strconv"
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
	arrImageID []int,
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
	if len(arrImageID) > 0 {
		q := "id IN (?" + strings.Repeat(",?", len(arrImageID)-1) + ")"
		wheres = append(wheres, q)
		for _, id := range arrImageID {
			values = append(values, strconv.Itoa(id))
		}
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

func DeleteImages(images []Image) error {
	tx, err := dbMap.Begin()
	if err != nil {
		return err
	}

	for _, i := range images {
		if _, err = tx.Delete(&i); err != nil {
			tx.Rollback()
			return err
		}

		path := strings.Replace(i.Path, "./img", imgRoot, 1)
		thumbPath := strings.Replace(path, "/img", "/img/thumb", 1)
		log.Println(path)
		log.Println(thumbPath)

		// Check file is exists
		if fileExists(path) {
			if _, err = exec.Command("rm", path).Output(); err != nil {
				log.Printf("rm error: %s", path)
			}
		}
		if fileExists(thumbPath) {
			if _, err = exec.Command("rm", thumbPath).Output(); err != nil {
				log.Printf("rm error: %s", thumbPath)
			}
		}
	}

	return tx.Commit()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (is *Images) Interface() []interface{} {
	imagesInterface := make([]interface{}, len(*is))
	for i, v := range *is {
		imagesInterface[i] = v
	}

	return imagesInterface
}
