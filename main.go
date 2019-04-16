package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/apple_watch_3"
	"log"
	"net/http"
)

func main() {
	initialiseRouter(":80")
}

func initialiseRouter(address string) {
	// Make a router and provide basic endpoints
	r := mux.NewRouter()
	r.HandleFunc("/health", health)
	r.HandleFunc("/", helloWorld)

	// Initialise the different data endpoints
	apple_watch_3.MakeStandardAppleWatch3Handler(r)

	// Start server. Log fatal if it crashes.
	log.Fatal(http.ListenAndServe(address, r))
}

func health(w http.ResponseWriter, _ *http.Request) {
	// Health endpoint for load balancer health checks
	w.WriteHeader(http.StatusOK)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello")
}
