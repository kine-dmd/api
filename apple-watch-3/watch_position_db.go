package apple_watch_3

import (
	"github.com/kine-dmd/api/dynamoDB"
	"log"
	"time"
)

type watchPositionDB interface {
	getWatchPosition(uuid string) (watchPosition, bool)
}

// Need a time interface for mocking time in tests
type apiTime interface {
	currentTime() time.Time
}

type systemTime struct{}

func (systemTime) currentTime() time.Time {
	return time.Now()
}

type dynamoCachedWatchDB struct {
	dbConnection  dynamoDB.DynamoDBInterface
	cache         map[string]watchPosition
	lastUpdatedAt time.Time
	timeKeeper    apiTime
}

func makeDynamoCachedWatchDB() *dynamoCachedWatchDB {
	dcw := new(dynamoCachedWatchDB)

	// Connect to Dynamo
	const table_name string = "apple_watch_3_positions"
	dcw.dbConnection = &dynamoDB.DynamoDBClient{}
	err := dcw.dbConnection.InitConn(table_name)
	if err != nil {
		log.Println("Error establishing connection to DynamoDB")
		log.Fatal(err)
	}

	// Use the standard clock and eager load the cache
	dcw.timeKeeper = systemTime{}
	dcw.updateCache()

	return dcw
}

func (dcw *dynamoCachedWatchDB) getWatchPosition(uuid string) (watchPosition, bool) {
	// If it had been more than 2 hours, always update the cache
	durationSinceUpdate := dcw.timeKeeper.currentTime().Sub(dcw.lastUpdatedAt)
	if durationSinceUpdate.Hours() >= 2 {
		dcw.updateCache()
	}

	// Try and retrieve the item from the updated cache
	val, ok := dcw.cache[uuid]
	if ok {
		return val, ok
	}

	// If item doesn't exist and cache was updated more than 15 minutes ago, retry update
	if durationSinceUpdate.Minutes() > 15 {
		dcw.updateCache()
	}

	// Return value whether it exists or not
	val, ok = dcw.cache[uuid]
	return val, ok
}

func (dcw *dynamoCachedWatchDB) updateCache() {
	// Get the raw data from DynamoDB
	unparsedRows := dcw.dbConnection.GetTableScan()

	// Format the data
	var parsedRows = make(map[string]watchPosition)
	for _, row := range unparsedRows {
		parsedRows[row["uuid"].(string)] = watchPosition{
			row["patientId"].(string),
			uint8(row["limb"].(float64)),
		}
	}

	// Update the cached values
	dcw.cache = parsedRows
	dcw.lastUpdatedAt = dcw.timeKeeper.currentTime()
}
