package watch_position_db

import (
	"github.com/kine-dmd/api/dynamoDB"
	"log"
)

type dynamoWatchDatabase struct {
	dbConnection dynamoDB.DynamoDBInterface
}

func makeStandardDynamoWatchDatabase() *dynamoWatchDatabase {
	// Connect to Dynamo
	const table_name string = "apple_watch_3_positions"
	dbConnection := &dynamoDB.DynamoDBClient{}
	err := dbConnection.InitConn(table_name)
	if err != nil {
		log.Println("Error establishing connection to DynamoDB")
		log.Fatal(err)
	}

	return makeDynamoWatchDatabase(dbConnection)
}

func makeDynamoWatchDatabase(dbInterface dynamoDB.DynamoDBInterface) *dynamoWatchDatabase {
	dwb := new(dynamoWatchDatabase)
	dwb.dbConnection = dbInterface
	return dwb
}

func (dwb dynamoWatchDatabase) GetTableScan() map[string]WatchPosition {
	// Get the raw data from DynamoDB
	unparsedRows := dwb.dbConnection.GetTableScan()

	// Format the data
	var parsedRows = make(map[string]WatchPosition)
	for _, row := range unparsedRows {
		parsedRows[row["uuid"].(string)] = WatchPosition{
			row["patientId"].(string),
			uint8(row["limb"].(float64)),
		}
	}

	return parsedRows
}

func (dwb dynamoWatchDatabase) GetWatchPosition(uuid string) (WatchPosition, bool) {
	table := dwb.GetTableScan()
	val, ok := table[uuid]
	return val, ok
}
