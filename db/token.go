package db

import (
	"aista-search/config"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
)

type Token struct {
	ID        int       `db:"id" json:"id"`
	token     string    `db:"token" json:"token"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

var salt string

func getSalt() string {
	if salt == "" {
		salt = config.GetEnv("DB_TOKEN_SALT", "episode solo")
	}

	return salt
}

func IsExistsToken() (bool, error) {
	query := "SELECT COUNT(*) FROM tokens"
	cnt, err := dbMap.SelectInt(query)
	checkErr(err)

	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

func CreateToken(md5 string) error {
	token := makeHash(md5)

	query := "INSERT INTO tokens (token) VALUES (?)"
	if _, err := dbMap.Exec(query, token); err != nil {
		return err
	}

	return nil
}

func IsValidToken(t string) bool {
	token := makeHash(t)

	query := "SELECT COUNT(*) FROM tokens WHERE token=?"
	cnt, err := dbMap.SelectInt(query, token)
	if err != nil {
		log.Println(err)
		return false
	}

	if cnt != 1 {
		log.Println(cnt)
		return false
	}

	return true
}

func makeHash(token string) string {
	hashByte := sha256.Sum256([]byte(token + getSalt()))
	return hex.EncodeToString(hashByte[:])
}
