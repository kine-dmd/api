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
			"limb":      float64(limbs[i]),
		}
	}

	return allRows
}
