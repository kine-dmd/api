package apple_watch_3

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/watch_position_db"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type apple_watch_3_handler struct {
	dataWriter Aw3DataWriter
	watchDB    watch_position_db.WatchPositionDatabase
}

func MakeStandardAppleWatch3Handler(r *mux.Router) *apple_watch_3_handler {
	// Always use a cached watch database
	watchDB := watch_position_db.MakeStandardCachedWatchDB()

	// Choose which data writer to use depending on local or AWS run
	if os.Getenv("kine_dmd_api_location") == "local" {
		return MakeAppleWatch3Handler(r, MakeStandardLocalFileDataWriter(), watchDB)
	}
	return MakeAppleWatch3Handler(r, MakeStandardKinesisDataWriter(), watchDB)
}

func MakeAppleWatch3Handler(r *mux.Router, queue Aw3DataWriter, watchDB watch_position_db.WatchPositionDatabase) *apple_watch_3_handler {
	// Assign the databases
	aw3Handler := new(apple_watch_3_handler)
	aw3Handler.dataWriter = queue
	aw3Handler.watchDB = watchDB

	// Pick a URL to handle
	r.HandleFunc("/upload/apple-watch-3/{uuid}", aw3Handler.binaryHandler).Methods("POST")
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

	// Check we have an integer number of rows
	if len(data)%ROW_SIZE_BYTES != 0 {
		http.Error(writer, "Non integer number of rows", http.StatusBadRequest)
		return
	}

	// Get the watch position
	position, exists := aw3Handler.watchDB.GetWatchPosition(watchId)
	if !exists {
		http.Error(writer, "Unable to match identifier.", http.StatusBadRequest)
		return
	}

	// Package the binary body together along with the watchId
	structuredData := UnparsedAppleWatch3Data{WatchPosition: position, RawData: data}

	// Send it to the relevant kinesis dataWriter
	err = aw3Handler.dataWriter.writeData(structuredData)
	if err != nil {
		http.Error(writer, "Server unable to forward body", http.StatusInternalServerError)
	}

	// Return the file number (if it exists else return empty)
	fileNum := request.Header.Get("Content-Disposition")
	_, _ = fmt.Fprintf(writer, fileNum)
}

func isValidUUID(uid string) bool {
	// Try and parse the uuid to a string
	_, err := uuid.FromString(uid)
	if err != nil {
		return false
	}
	return true
}
