package db

import (
	"aista-search/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type ValidateError struct {
	msg string
}

func (err *ValidateError) Error() string {
	return err.msg
}

func newValidateError(s string) *ValidateError {
	err := new(ValidateError)
	err.msg = s
	return err
}

var dbMap *gorp.DbMap

func Get() *gorp.DbMap {
	return dbMap
}

func Connect() {
	driver := config.GetEnv("DB_DRIVER", "mysql")
	url := config.GetEnv("DB_URL", "")
	db, err := sql.Open(driver, url)
	if err != nil {
		panic(err)
	}

	dbMap = initDb(db)
}

func initDb(db *sql.DB) *gorp.DbMap {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Episode{}, "episodes")
	dbmap.AddTableWithName(Image{}, "images").SetKeys(true, "ID")
	dbmap.AddTableWithName(Token{}, "tokens").SetKeys(true, "ID")
	dbmap.CreateTablesIfNotExists()

	return dbmap
}
