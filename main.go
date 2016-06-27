package main

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"aista-search/session"
	"aista-search/view"
	"net/http"
)

func main() {
	config.LoadEnv()
	db.Connect()
	session.Configure()
	view.Configure()

	router := route.New()
	http.ListenAndServe(":8080", router)
}
