package apple_watch_3

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

// A router to use when testing
var router *mux.Router

// Run before each test
func init() {
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

func TestNilBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", nil)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestNilBodyWithSetContentLength(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", nil)
	req.ContentLength = 55
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func Test0LenBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func Test0LenBodyWithSetContentLength(t *testing.T) {
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
	req.ContentLength = 55
	response := sendRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestValidUUID(t *testing.T) {
	checkValidUUID(t, "00000000-0000-0000-0000-000000000000")
}

func TestValidUUIDWithoutDashes(t *testing.T) {
	checkValidUUID(t, "00000000000000000000000000000000")
}

func checkValidUUID(t *testing.T, validUUID string) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDoer := mocks.NewMockKinesisQueueInterface(mockCtrl)
	queue = mockDoer

	// Exactly one thing should be sent to the queue. Partition key should be UUID
	mockDoer.EXPECT().SendToQueue(gomock.Any(), validUUID).Return(nil).Times(1)

	// Make and send a request with some data
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/"+validUUID, body)
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
