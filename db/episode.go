package db

import (
	"time"
)

type Episode struct {
	Id        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Status    uint      `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

const (
	Created = iota
	Collected
	Crassified
	Extracted
	Registered
)

func (e *Episode) GetStatus() string {
	switch e.Status {
	case Created:
		return "作成済み"
	case Collected:
		return "画像収集済み"
	case Crassified:
		return "振り分け済み"
	case Extracted:
		return "文字抽出済み"
	case Registered:
		return "登録済み"
	default:
		return "未分類"
	}

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

func UpdateEpisode(id int, title string, status uint) error {
	var query string
	if title == "" {
		query = "UPDATE episodes SET status=? WHERE id=?"
		if _, err := dbMap.Exec(query, status, id); err != nil {
			return err
		}
	} else {
		query = "UPDATE episodes SET title=?, status=? WHERE id=?"
		if _, err := dbMap.Exec(query, title, status, id); err != nil {
			return err
		}
	}

	return nil
}

func GetLatestEpisode() (*Episode, error) {
	var ep Episode
	query := "SELECT * FROM episodes ORDER BY id DESC LIMIT 1"
	if err := dbMap.SelectOne(&ep, query); err != nil {
		return nil, err
	}

	return &ep, nil
}

// unique check
func (e *Episode) validate() (bool, error) {
	count, _ := dbMap.SelectInt("SELECT count(*) FROM episodes WHERE id=?", e.Id)
	if count > 0 {
		return false, newValidateError("Duplicate record in episodes")
	}

	return true, nil
}
