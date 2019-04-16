package main

import (
	"net/http"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	// Start server locally on new thread
	go initialiseRouter(":8888")

	// Allow 2 seconds for the server to start before sending request
	time.Sleep(2 * time.Second)
	resp, err := http.Get("http://localhost:8888/health")

	// Check no error was received from the endpoint
	if err != nil {
		t.Fatalf("Got error from health endpoint: %s", err)
	}

	// Check the status code is 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Received status code %d from health endpoint", resp.StatusCode)
	}
}
