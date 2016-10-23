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

var (
	dbMap   *gorp.DbMap
	imgRoot string
)

func Get() *gorp.DbMap {
	return dbMap
}

func Connect() {
	driver := config.GetEnv("DB_DRIVER", "mysql")

	// Connect mysql
	db, err := sql.Open(driver, config.GetEnv("DB_BASE_URL", ""))
	if err != nil {
		panic(err)
	}

	initDb(db)
	// Connect database
	db, err = sql.Open(driver, config.GetEnv("DB_URL", ""))
	if err != nil {
		panic(err)
	}
	dbMap = initTable(db)
	imgRoot = config.GetEnv("IMG_ROOT", "./img")
}

func initDb(db *sql.DB) {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	// create db if not exists
	dbName := config.GetEnv("DB_NAME", "aista_search_dev")
	sql := "CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8;"
	if _, err := dbmap.Exec(sql); err != nil {
		panic(err)
	}
}

func initTable(db *sql.DB) *gorp.DbMap {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Episode{}, "episodes")
	dbmap.AddTableWithName(Image{}, "images").SetKeys(true, "ID")
	dbmap.AddTableWithName(Token{}, "tokens").SetKeys(true, "ID")
	dbmap.CreateTablesIfNotExists()

	return dbmap
}
