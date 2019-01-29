package apple_watch_3

import (
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/kinesisqueue"
	"io/ioutil"
	"log"
	"net/http"
)

type appleWatch3Data struct {
	UUID string `json:"UUID"`
	Data []byte `json:"Data"`
}

const STREAM_NAME = "apple-watch-3"

var queue kinesisqueue.KinesisQueueInterface = &kinesisqueue.KinesisQueueClient{}

func Init(r *mux.Router) {
	// Open a kinesis queue connection
	err := queue.InitConn(STREAM_NAME)
	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/upload/apple-watch-3/{uuid}", binaryHandler).Methods("POST")
}

func binaryHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the UUID from the URL
	vars := mux.Vars(request)
	uuid := vars["uuid"]

	// Check there is a body
	if request.ContentLength == 0 {
		http.Error(writer, "0 length body.", http.StatusExpectationFailed)
		return
	}

	// Entire body is data so read
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Unable to read body from Apple Watch 3. UUID: ", uuid)
		http.Error(writer, "Unable to read POST body data", http.StatusUnprocessableEntity)
		return
	}

	// Package the binary data together along with the uuid
	structuredData := appleWatch3Data{UUID: uuid, Data: data}

	// Send it to the relevant kinesis queue
	err = queue.SendToQueue(structuredData, uuid)
	if err != nil {
		http.Error(writer, "Server unable to forward data", http.StatusInternalServerError)
	}
}
