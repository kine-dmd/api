package apple_watch_3

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

// A router to use when testing
var router *mux.Router

// Run before each test
func init() {
	print("Init run")
	router = mux.NewRouter()
	Init(router)
}

func TestWrongLengthUUID(t *testing.T) {
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/someuuid", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestUUIDNotBase64(t *testing.T) {
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/ZZZZZZZZ-XXXX-YYYY-UUUU-WWWWWWWWWWWW", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestValidUUID(t *testing.T) {
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestValidUUIDWithoutDashes(t *testing.T) {
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000000000000000000000000001", body)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestNilBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", nil)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func Test0LenBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
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
