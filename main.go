package main

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"aista-search/session"
	"aista-search/view"
	"aista-search/view/plugin"
	"net/http"
)

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
	http.ListenAndServe(":8080", router)
}
