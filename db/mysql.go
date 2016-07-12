package db

import (
	"aista-search/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

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
	dbmap.CreateTablesIfNotExists()

	return dbmap
}
