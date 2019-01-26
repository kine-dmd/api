package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", health)
	r.HandleFunc("/", helloWorld)

	log.Fatal(http.ListenAndServe(":80", r))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello")
}
