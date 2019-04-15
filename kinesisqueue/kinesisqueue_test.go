package kinesisqueue

import "testing"

func TestSendingToQueue(t *testing.T) {
	// Set up a connection to the test stream
	kq := MakeKinesisQueue("testStream")

	// Send a basic payload
	payload := map[string]int{"key1": 1}
	err := kq.SendToQueue(payload, "key")
	if err != nil {
		t.Fatal(err)
	}
}
