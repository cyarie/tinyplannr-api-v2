package main

import (
	"log"
	"net/http"

	"github.com/cyarie/tinyplannr-api-v2/api"
)

func main() {
	// Simple! Just pull in the router from the API package, and let 'er roll.
	router := api.Api()
	log.Println(http.ListenAndServe(":8080", router))
}
