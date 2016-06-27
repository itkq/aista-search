package main

import (
	"aista-search/route"
	"net/http"
)

func main() {
	router := route.New()

	http.ListenAndServe(":8080", router)
}
