package watch_position_db

import (
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/dynamoDB"
	"testing"
)

func TestScanDynamoDatabase(t *testing.T) {
	// Exactly 0 things should be sent to the queue
	mockCtrl, mockDB := makeMockDB(t)
	mockDB.EXPECT().GetTableScan().Times(1).Return(makeTestRows()).Times(1)

	// Make a watch position database using the raw mocked dynamo instance and query it
	watchDB := makeDynamoWatchDatabase(mockDB)
	scan := watchDB.GetTableScan()

	compareResults(scan, t)
	mockCtrl.Finish()
}

func makeMockDB(t *testing.T) (*gomock.Controller, *dynamoDB.MockDynamoDBInterface) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockDB := dynamoDB.NewMockDynamoDBInterface(mockCtrl)
	return mockCtrl, mockDB
}

func makeTestRows() []map[string]interface{} {
	uuids, patiendIds, limbs := makeFakePositionData()

	// Match the fake data to the format given by DynamoDB
	allRows := make([]map[string]interface{}, 3)
	for i := range allRows {
		allRows[i] = map[string]interface{}{
			"uuid":      uuids[i],
			"patientId": patiendIds[i],
			"limb":      limbs[i],
		}
	}

	return allRows
}

func makeFakePositionData() ([]string, []string, []float64) {
	// Fake data in use
	uuids := []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}
	patiendIds := []string{"dmd01", "dmd02", "dmd03"}
	limbs := []float64{1, 2, 3}
	return uuids, patiendIds, limbs
}

func compareResults(scan map[string]WatchPosition, t *testing.T) {
	origUUIDs, origPatientIds, origLimbs := makeFakePositionData()

	// Check each row in the original data still exists
	for i := 0; i < len(origUUIDs); i++ {
		val, ok := scan[origUUIDs[i]]

		// Check that the 3 values all exist and match up to their original values
		if !ok {
			t.Errorf("UUID %s not found when should've been", origUUIDs[i])
		}
		if val.PatientID != origPatientIds[i] {
			t.Errorf("Mismatching patient identifiers. Expected %s . Got %s .", origPatientIds[i], val.PatientID)
		}
		if val.Limb != uint8(origLimbs[i]) {
			t.Errorf("Mismatching limb identifiers. Expected %f . Got %d .", origLimbs[i], val.Limb)
		}

	}
}
