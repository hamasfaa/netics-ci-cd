package main

import (
	"learn-ci-cd/handler"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	timeZone := time.FixedZone("UTC+7", 8*60*60)
	timeUp := time.Now().In(timeZone).Format("2006-01-02 15:04:05")

	route := mux.NewRouter()

	// http.HandleFunc("/", handler.IndexHandler())
	// http.HandleFunc("/health", handler.HealthHandler(timeUp, timeZone))

	route.HandleFunc("/", handler.IndexHandler())
	route.HandleFunc("/testing", handler.TestingHandler())
	route.HandleFunc("/health", handler.HealthHandler(timeUp, timeZone))

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", route); err != nil {
		log.Fatal(err)
	}
}
