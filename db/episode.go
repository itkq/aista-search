package db

import (
	"time"
)

type Episode struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Status    uint      `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

const (
	EpCreated = iota + 1
	EpRetrieved
	EpRegistered
)

func (e *Episode) GetStatus() string {
	switch e.Status {
	case EpCreated:
		return "作成済み"
	case EpRetrieved:
		return "画像収集済み"
	case EpRegistered:
		return "登録済み"
	default:
		return "未分類"
	}
}

func GetEpisode(id int) (*Episode, error) {
	var ep Episode
	query := "SELECT * FROM episodes WHERE id=?"
	if err := dbMap.SelectOne(&ep, query, id); err != nil {
		return nil, err
	}

	return &ep, nil
}

func GetEpisodes(cnt int) (*[]Episode, error) {
	var episodes []Episode
	query := "SELECT * FROM episodes ORDER BY id LIMIT ?"
	if _, err := dbMap.Select(&episodes, query, cnt); err != nil {
		return nil, err
	}

	return &episodes, nil
}

func CreateEpisode(id int, title string, status uint) (*Episode, error) {
	episode := Episode{
		ID:        id,
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
	var values []interface{}

	if title == "" {
		query = "UPDATE episodes SET status=? WHERE id=?"
		values = append(values, interface{}(status))
		values = append(values, interface{}(id))
	} else {
		query = "UPDATE episodes SET title=?, status=? WHERE id=?"
		values = append(values, interface{}(title))
		values = append(values, interface{}(status))
		values = append(values, interface{}(id))
	}
	if _, err := dbMap.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

// unique check
func (e *Episode) validate() (bool, error) {
	count, _ := dbMap.SelectInt("SELECT count(*) FROM episodes WHERE id=?", e.ID)
	if count > 0 {
		return false, newValidateError("Duplicate record in episodes")
	}

	return true, nil
}
