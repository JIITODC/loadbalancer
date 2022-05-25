package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// get port and service no. from env
	port := os.Getenv("PORT")
	serviceNo := os.Getenv("SERVICE")

	// set defaults
	if port == "" {
		port = "8080"
	}

	// create new logger
	logger := log.New(os.Stdout, fmt.Sprintf("Service %s: ", serviceNo), 2)

	// handlers
	myMux := http.NewServeMux()
	myMux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		return
	})
	myMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			logger.Printf("%s request received...\n", r.Method)
		}
	})

	// start serving
	logger.Printf("Listening on port %s...", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), myMux)
}
