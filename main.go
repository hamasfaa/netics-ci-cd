package main

import (
	"learn-ci-cd/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.IndexHandler())
	// http.HandleFunc("/health", handler.HealthHandler())
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
