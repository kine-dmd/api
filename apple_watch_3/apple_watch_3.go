package apple_watch_3

import (
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/kinesisqueue"
	"github.com/kine-dmd/api/watch_position_db"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
)

type unparsedAppleWatch3Data struct {
	WatchPosition watch_position_db.WatchPosition `json:"WatchPosition"`
	RawData       []byte                          `json:"RawData"`
}

type apple_watch_3_handler struct {
	queue   kinesisqueue.KinesisQueueInterface
	watchDB watch_position_db.WatchPositionDB
}

func MakeStandardAppleWatch3Handler(r *mux.Router) *apple_watch_3_handler {
	// Open a kinesis queue & dynamo DB connection
	const STREAM_NAME = "apple_watch_3"
	queue := kinesisqueue.MakeKinesisQueue(STREAM_NAME)
	watchDB := watch_position_db.MakeStandardDynamoCachedWatchDB()

	return MakeAppleWatch3Handler(r, queue, watchDB)
}

func MakeAppleWatch3Handler(r *mux.Router, queue kinesisqueue.KinesisQueueInterface, watchDB watch_position_db.WatchPositionDB) *apple_watch_3_handler {
	// Assign the databases
	aw3Handler := new(apple_watch_3_handler)
	aw3Handler.queue = queue
	aw3Handler.watchDB = watchDB

	// Pick a URL to handle
	r.HandleFunc("/upload/apple_watch_3/{uuid}", aw3Handler.binaryHandler).Methods("POST")
	return aw3Handler
}

func (aw3Handler apple_watch_3_handler) binaryHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the UUID from the URL
	vars := mux.Vars(request)
	watchId := vars["uuid"]

	// Check the uuid is valid
	if !isValidUUID(watchId) {
		http.Error(writer, "Bad UUID.", http.StatusBadRequest)
		return
	}

	// Check there is a body
	if request.ContentLength == 0 || request.Body == nil {
		http.Error(writer, "0 length body.", http.StatusBadRequest)
		return
	}

	// Entire body is body so read
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Unable to read body from Apple Watch 3. UUID: ", watchId)
		http.Error(writer, "Unable to read POST body body", http.StatusBadRequest)
		return
	}
	if data == nil || len(data) == 0 {
		http.Error(writer, "0 length body.", http.StatusBadRequest)
		return
	}

	// Get the watch position
	position, exists := aw3Handler.watchDB.GetWatchPosition(watchId)
	if !exists {
		http.Error(writer, "Unable to match identifier.", http.StatusBadRequest)
		return
	}

	// Package the binary body together along with the watchId
	structuredData := unparsedAppleWatch3Data{WatchPosition: position, RawData: data}

	// Send it to the relevant kinesis queue
	err = aw3Handler.queue.SendToQueue(structuredData, watchId)
	if err != nil {
		http.Error(writer, "Server unable to forward body", http.StatusInternalServerError)
	}
}

func isValidUUID(uid string) bool {
	// Try and parse the uuid to a string
	_, err := uuid.FromString(uid)
	if err != nil {
		return false
	}
	return true
}