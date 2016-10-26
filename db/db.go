package db

import (
	"aista-search/config"
	"aista-search/util"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"os"
	"os/exec"
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
	baseUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/",
		config.GetEnv("DB_USER", "root"),
		config.GetEnv("DB_PASS", ""),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "3306"),
	)
	db, err := sql.Open(driver, baseUrl)
	checkErr(err)

	initDb(db)
	// Connect database
	dbUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s%s",
		config.GetEnv("DB_USER", "root"),
		config.GetEnv("DB_PASS", ""),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "3306"),
		config.GetEnv("DB_NAME", "aista_search_dev"),
		config.GetEnv("DB_OPT", ""),
	)
	db, err = sql.Open(driver, dbUrl)
	checkErr(err)

	dbMap = initTable(db)
	imgRoot = config.GetEnv("IMG_ROOT", "./img")

	tf, err := IsExistsToken()
	checkErr(err)
	if !tf {
		key := util.RandomStr(16)
		err = CreateToken(key)
		checkErr(err)

		envFile := fmt.Sprintf(".env.%s", os.Getenv("GO_ENV"))
		if config.GetEnv("API_TOKEN", "") != "" {
			sub := "s/API_TOKEN.*$/API_TOKEN=" + key + "/"
			exec.Command("sed", "-i", "-e", sub, envFile).Run()
		} else {
			fp, err := os.OpenFile(envFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			checkErr(err)

			defer fp.Close()
			writer := bufio.NewWriter(fp)
			writer.WriteString("API_TOKEN=" + key + "\n")
			writer.Flush()
		}
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
