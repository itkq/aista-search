package db

import (
	"time"
)

type Episode struct {
	Id        int       `db:"id"`
	Title     string    `db:"title"`
	Status    uint      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
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

	if _, err := episode.validate(); err != nil {
		return nil, err
	}

	if err := dbMap.Insert(&episode); err != nil {
		return nil, err
	}

	return &episode, nil
}

// unique check
func (e *Episode) validate() (bool, error) {
	count, _ := dbMap.SelectInt("SELECT count(*) FROM episodes WHERE id=?", e.Id)
	if count > 0 {
		return false, newValidateError("Duplicate record in episodes")
	}

	return true, nil
}
