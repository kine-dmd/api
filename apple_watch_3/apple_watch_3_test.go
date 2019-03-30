package apple_watch_3

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/kine-dmd/api/kinesisqueue"
	"github.com/kine-dmd/api/mocks"
	"github.com/kine-dmd/api/watch_position_db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrongLengthUUID(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().SendToQueue(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/someuuid", body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestUUIDNotBase64(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().SendToQueue(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/ZZZZZZZZ-XXXX-YYYY-UUUU-WWWWWWWWWWWW", body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestNilBody(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().SendToQueue(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/00000000-0000-0000-0000-000000000000", nil)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func TestNilBodyWithSetContentLength(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().SendToQueue(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/00000000-0000-0000-0000-000000000000", nil)
	req.ContentLength = 55
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func Test0LenBody(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().SendToQueue(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	mockCtrl.Finish()
}

func Test0LenBodyWithSetContentLength(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockQueue, mockDB := makeMockQueueAndDB(t)
	router, _ := initRouterAndHandler(mockQueue, mockDB)
	mockQueue.EXPECT().SendToQueue(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	// Set a mismatching content length
	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/00000000-0000-0000-0000-000000000000", bytes.NewReader([]byte{}))
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
	mockQueue.EXPECT().SendToQueue(gomock.Any(), validUUID).Return(nil).Times(1)
	mockDB.EXPECT().GetWatchPosition(validUUID).Return(watch_position_db.WatchPosition{"dmd01", 1}, true).Times(1)

	// Make and send a request with some data
	body := bytes.NewReader([]byte{1, 2, 3})
	req, _ := http.NewRequest("POST", "/upload/apple_watch_3/"+validUUID, body)
	response := sendRequest(router, req)
	checkResponseCode(t, http.StatusOK, response.Code)

	mockCtrl.Finish()
}

func makeMockQueueAndDB(t *testing.T) (*gomock.Controller, *mocks.MockKinesisQueueInterface, *mocks.MockWatchPositionDB) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockQueue := mocks.NewMockKinesisQueueInterface(mockCtrl)
	mockDB := mocks.NewMockWatchPositionDB(mockCtrl)
	return mockCtrl, mockQueue, mockDB
}

func initRouterAndHandler(mockQueue kinesisqueue.KinesisQueueInterface, mockWatchDB watch_position_db.WatchPositionDB) (*mux.Router, *apple_watch_3_handler) {
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
