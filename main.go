package main

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"aista-search/session"
	"aista-search/view"
	"aista-search/view/plugin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	config.LoadEnv()
	db.Connect()
	session.Configure()

	view.Configure()
	view.LoadPlugins(
		plugin.FormattedTime(),
		plugin.EpisodeStatus(),
	)

	router := route.New()
	http.ListenAndServe(config.GetEnv("APP_PORT", ":8080"), router)
}
