package watch_position_db

import (
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/mocks/mock_time"
	"sync"
	"testing"
	"time"
)

func TestGetsDataOnCreation(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	mockTime.EXPECT().CurrentTime().Return(time.Now()).Times(2)

	// Make an empty cached DB and query it
	_ = MakeCachedWatchDB(mockDB, mockTime)
	mockCtrl.Finish()
}

func TestRetrievingRowFromCache(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	curTime := time.Now()
	mockTime.EXPECT().CurrentTime().Return(curTime).Times(2)

	// Make an empty cached DB and query it
	dcw := MakeCachedWatchDB(mockDB, mockTime)
	mockTime.EXPECT().CurrentTime().Return(curTime.Add(time.Hour)).Times(1)
	watchPos, exists := dcw.GetWatchPosition("00000000-0000-0000-0000-000000000001")

	// Check that no result was obtained
	checkRetrievedExistence(t, exists, true)
	checkRetrievedWatchPositionValues(t, watchPos, 2, "dmd01")
	mockCtrl.Finish()
}

func TestRetrievingRowReloadCache(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)

	// Make an empty cached DB and query it
	curTime := time.Now()
	mockTime.EXPECT().CurrentTime().Return(curTime).Times(2)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	dcw := MakeCachedWatchDB(mockDB, mockTime)

	// Query the data 3 hours after load
	mockTime.EXPECT().CurrentTime().Return(curTime.Add(time.Hour * 3)).Times(3)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	watchPos, exists := dcw.GetWatchPosition("00000000-0000-0000-0000-000000000001")

	// Check that no result was obtained
	checkRetrievedExistence(t, exists, true)
	checkRetrievedWatchPositionValues(t, watchPos, 2, "dmd01")
	mockCtrl.Finish()
}

func TestRetrievingRowReloadCacheThenFromCache(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)

	// Make an empty cached DB and query it
	curTime := time.Now()
	mockTime.EXPECT().CurrentTime().Return(curTime).Times(2)
	dcw := MakeCachedWatchDB(mockDB, mockTime)

	// Query the data 3 hours after load
	mockTime.EXPECT().CurrentTime().Return(curTime.Add(time.Hour * 3)).Times(3)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	watchPos, exists := dcw.GetWatchPosition("00000000-0000-0000-0000-000000000001")

	// Check the results
	checkRetrievedExistence(t, exists, true)
	checkRetrievedWatchPositionValues(t, watchPos, 2, "dmd01")

	// Query the data 1 hours after previous load
	mockTime.EXPECT().CurrentTime().Return(curTime.Add(time.Hour * 4)).Times(1)
	watchPos, exists = dcw.GetWatchPosition("00000000-0000-0000-0000-000000000002")

	// Check the results
	checkRetrievedExistence(t, exists, true)
	checkRetrievedWatchPositionValues(t, watchPos, 1, "dmd02")
	mockCtrl.Finish()
}

func TestStressOnlyPerformsOneReload(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)

	// Make an empty cached DB and query it
	curTime := time.Now()
	mockTime.EXPECT().CurrentTime().Return(curTime).Times(2)
	dcw := MakeCachedWatchDB(mockDB, mockTime)

	// Query the data 3 hours after load
	mockTime.EXPECT().CurrentTime().Return(curTime.Add(time.Hour * 3)).AnyTimes()
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)

	// Simulate 50 requests being sent at once. Table scan should still only be called once
	waitGroup := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		waitGroup.Add(1)

		// Spin  up a new thread to send request
		go func() {
			// Make the request
			defer waitGroup.Done()
			watchPos, exists := dcw.GetWatchPosition("00000000-0000-0000-0000-000000000001")

			// Check the result was obtained
			checkRetrievedExistence(t, exists, true)
			checkRetrievedWatchPositionValues(t, watchPos, 2, "dmd01")
		}()
	}

	// Wait for all threads to finish requests and check mocks status
	waitGroup.Wait()
	mockCtrl.Finish()
}

func checkRetrievedExistence(t *testing.T, exists bool, expectedExists bool) {
	if exists != expectedExists {
		t.Errorf("Existence of watchposition incorrect. Expedcted %t got %t.", expectedExists, exists)
	}
}

func checkRetrievedWatchPositionValues(t *testing.T, position WatchPosition, limb uint8, patientID string) {
	if position.Limb != limb {
		t.Errorf("Mismatching watch positions. Expected %d got %d", limb, position.Limb)
	}

	if position.PatientID != patientID {
		t.Errorf("Mismatching patient id. Expected %s got %s", patientID, position.PatientID)
	}
}

func makeFakeData() map[string]WatchPosition {
	// Make a list of all rows
	allData := make(map[string]WatchPosition)

	// Create 3 fake rows
	allData["00000000-0000-0000-0000-000000000000"] = WatchPosition{"dmd01", 1}
	allData["00000000-0000-0000-0000-000000000001"] = WatchPosition{"dmd01", 2}
	allData["00000000-0000-0000-0000-000000000002"] = WatchPosition{"dmd02", 1}

	return allData
}

func makeMocks(t *testing.T) (*gomock.Controller, *MockWatchPositionDatabase, *mock_api_time.MockApiTime) {
	// Make a mock for the time and database
	mockCtrl := gomock.NewController(t)
	mockDB := NewMockWatchPositionDatabase(mockCtrl)
	mockTime := mock_api_time.NewMockApiTime(mockCtrl)
	return mockCtrl, mockDB, mockTime
}
