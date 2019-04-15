package apple_watch_3

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/watch_position_db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrongLengthUUID(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(0)

	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/someuuid", body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestUUIDNotBase64(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(0)

	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/ZZZZZZZZ-XXXX-YYYY-UUUU-WWWWWWWWWWWW", body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestNilBody(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", nil)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestNilBodyWithSetContentLength(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", nil)
	req.ContentLength = 55
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func Test0LenBody(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func Test0LenBodyWithSetContentLength(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(0)

	// Set a mismatching content length
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
	req.ContentLength = 55
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestValidUUID(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	validUUID := "00000000-0000-0000-0000-000000000000"
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().writeData(gomock.Any()).Return(nil).Times(1)
	mockDB.EXPECT().GetWatchPosition(validUUID).Return(watch_position_db.WatchPosition{"dmd01", 1}, true).Times(1)

	// Make and send a request with some data
	body := bytes.NewReader(makeValidByteInput(1))
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/"+validUUID, body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusOK, response.Code)

	mockCtrl.Finish()
}

func TestIncorrectLengthData(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	validUUID := "00000000-0000-0000-0000-000000000000"
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)

	// Make and send a request with some data
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/"+validUUID, body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestIncorrectLengthDataWithCorrectSetLength(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	validUUID := "00000000-0000-0000-0000-000000000000"
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)

	// Make and send a request with some data
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple-watch-3/"+validUUID, body)
	req.ContentLength = ROW_SIZE_BYTES
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func makeValidByteInput(numRows int) []byte {
	var watchData []byte
	for i := 0; i < numRows; i++ {
		watchData = append(watchData, make([]byte, 88)...)
	}
	return watchData
}

func makeMockQueueAndDB(t *testing.T) (*gomock.Controller, *MockAw3DataWriter, *watch_position_db.MockWatchPositionDatabase) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockQueue := NewMockAw3DataWriter(mockCtrl)
	mockDB := watch_position_db.NewMockWatchPositionDatabase(mockCtrl)
	return mockCtrl, mockQueue, mockDB
}

func initRouterAndHandler(mockQueue Aw3DataWriter, mockWatchDB watch_position_db.WatchPositionDatabase) (*mux.Router, *apple_watch_3_handler) {
	router := mux.NewRouter()
	handler := MakeAppleWatch3Handler(router, mockQueue, mockWatchDB)
	return router, handler
}

func sendRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
