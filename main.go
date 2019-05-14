package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/apple_watch_3"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Set up the URL router
	r := initialiseRouter()

	// Choose which server to run depending on local or AWS run
	if os.Getenv("kine_dmd_api_location") == "local" {
		log.Fatal(http.ListenAndServeTLS(":443", "/certs/apollo.cert", "/certs/apollo.key", r))
	}
	log.Fatal(http.ListenAndServe(":80", r))
}

func initialiseRouter() *mux.Router {
	// Make a router and provide basic endpoints
	r := mux.NewRouter()
	r.HandleFunc("/health", health)
	r.HandleFunc("/", helloWorld)

	// Initialise the different data endpoints
	apple_watch_3.MakeStandardAppleWatch3Handler(r)

	// Return the set up router
	return r
}

func health(w http.ResponseWriter, _ *http.Request) {
	// Health endpoint for load balancer health checks
	w.WriteHeader(http.StatusOK)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello. Time is %s", time.Now().String())
}
