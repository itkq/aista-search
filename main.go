package main

import (
	"aista-search/config"
	"aista-search/db"
	"aista-search/route"
	"net/http"
)

func main() {
	config.LoadEnv()
	db.Connect()

	router := route.New()
	http.ListenAndServe(":8080", router)
}
