package watch_position_db

import (
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/mocks/mock_dynamo_db"
	"github.com/kine-dmd/api/mocks/mock_time"
	"testing"
	"time"
)

func TestGetsDataOnCreation(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	mockTime.EXPECT().CurrentTime().Return(time.Now()).Times(2)

	// Make an empty cached DB and query it
	_ = MakeDynamoCachedWatchDB(mockDB, mockTime)
	mockCtrl.Finish()
}

func TestRetrievingRowFromCache(t *testing.T) {
	// Make mocks and set the expectations
	mockCtrl, mockDB, mockTime := makeMocks(t)
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)
	curTime := time.Now()
	mockTime.EXPECT().CurrentTime().Return(curTime).Times(2)

	// Make an empty cached DB and query it
	dcw := MakeDynamoCachedWatchDB(mockDB, mockTime)
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
	mockDB.EXPECT().GetTableScan().Return(makeFakeData()).Times(1)

	// Make an empty cached DB and query it
	curTime := time.Now()
	mockTime.EXPECT().CurrentTime().Return(curTime).Times(2)
	dcw := MakeDynamoCachedWatchDB(mockDB, mockTime)

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
	dcw := MakeDynamoCachedWatchDB(mockDB, mockTime)

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

func makeFakeDataRow(uuid string, limb uint8, patientId string) map[string]interface{} {
	return map[string]interface{}{"uuid": uuid,
		"limb":      float64(limb),
		"patientId": patientId,
	}
}

func makeFakeData() []map[string]interface{} {
	// Make a list of all rows
	allData := make([]map[string]interface{}, 3)

	// Create 3 fake rows
	allData[0] = makeFakeDataRow("00000000-0000-0000-0000-000000000000", 1, "dmd01")
	allData[1] = makeFakeDataRow("00000000-0000-0000-0000-000000000001", 2, "dmd01")
	allData[2] = makeFakeDataRow("00000000-0000-0000-0000-000000000002", 1, "dmd02")

	return allData
}

func makeMocks(t *testing.T) (*gomock.Controller, *mock_dynamoDB.MockDynamoDBInterface, *mock_api_time.MockApiTime) {
	// Make a mock for the time and database
	mockCtrl := gomock.NewController(t)
	mockDB := mock_dynamoDB.NewMockDynamoDBInterface(mockCtrl)
	mockTime := mock_api_time.NewMockApiTime(mockCtrl)
	return mockCtrl, mockDB, mockTime
}
