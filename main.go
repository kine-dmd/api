package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/apple_watch_3"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", health)
	r.HandleFunc("/", helloWorld)

	// Initialise the different endpoint adapters
	apple_watch_3.MakeStandardAppleWatch3Handler(r)

	// Start server. Log fatal if it crashes.
	log.Fatal(http.ListenAndServe(":80", r))
}

func health(w http.ResponseWriter, r *http.Request) {
	// Health endpoint for load balancer health checks
	w.WriteHeader(http.StatusOK)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello")
}
