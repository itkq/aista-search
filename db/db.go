package db

import (
	"aista-search/config"
	"aista-search/util"
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"os"
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
	var err error

	driver := config.GetEnv("DB_DRIVER", "mysql")

	// Connect mysql
	db, err := sql.Open(driver, config.GetEnv("DB_BASE_URL", ""))
	checkErr(err)

	initDb(db)
	// Connect database
	db, err = sql.Open(driver, config.GetEnv("DB_URL", ""))
	checkErr(err)

	dbMap = initTable(db)
	imgRoot = config.GetEnv("IMG_ROOT", "./img")

	tf, err := IsExistsToken()
	log.Println(tf)
	checkErr(err)
	if !tf {
		key := util.RandomStr(16)
		err = CreateToken(key)
		checkErr(err)

		envFile := config.GetEnv("ENV_FILE", ".envrc")
		fp, err := os.OpenFile(envFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		checkErr(err)

		defer fp.Close()
		writer := bufio.NewWriter(fp)
		writer.WriteString("export API_TOKEN=" + key + "\n")
		writer.Flush()
	}
}

func initDb(db *sql.DB) {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	// create db if not exists
	dbName := config.GetEnv("DB_NAME", "aista_search_dev")
	sql := "CREATE DATABASE IF NOT EXISTS " + dbName + " DEFAULT CHARACTER SET utf8;"
	_, err := dbmap.Exec(sql)
	checkErr(err)
}

func initTable(db *sql.DB) *gorp.DbMap {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Episode{}, "episodes")
	dbmap.AddTableWithName(Image{}, "images").SetKeys(true, "ID")
	dbmap.AddTableWithName(Token{}, "tokens").SetKeys(true, "ID")
	dbmap.CreateTablesIfNotExists()

	return dbmap
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
