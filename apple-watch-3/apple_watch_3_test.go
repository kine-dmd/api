package apple_watch_3

import (
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// A router to use when testing
var router *mux.Router = mux.NewRouter()
var body io.Reader = bytes.NewReader([]byte{0, 0, 0})

// Run before each test
func init() {
	Init(router)
}

func TestWrongLengthUUID(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/someuuid", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestUUIDNotBase64(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/ZZZZZZZZ-XXXX-YYYY-UUUU-WWWWWWWWWWWW", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestValidUUID(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestValidUUIDWithoutDashes(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000000000000000000000000000", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func sendRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
