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
type ApiTime interface {
	CurrentTime() time.Time
}

type systemTime struct{}

func (systemTime) CurrentTime() time.Time {
	return time.Now()
}

type dynamoCachedWatchDB struct {
	dbConnection  dynamoDB.DynamoDBInterface
	cache         map[string]watchPosition
	lastUpdatedAt time.Time
	timeKeeper    ApiTime
}

func makeStandardDynamoCachedWatchDB() *dynamoCachedWatchDB {
	// Connect to Dynamo
	const table_name string = "apple_watch_3_positions"
	dbConnection := &dynamoDB.DynamoDBClient{}
	err := dbConnection.InitConn(table_name)
	if err != nil {
		log.Println("Error establishing connection to DynamoDB")
		log.Fatal(err)
	}

	return makeDynamoCachedWatchDB(dbConnection, systemTime{})
}

func makeDynamoCachedWatchDB(dbConnection dynamoDB.DynamoDBInterface, timeKeeper ApiTime) *dynamoCachedWatchDB {
	// Create a new connection
	dcw := new(dynamoCachedWatchDB)
	dcw.dbConnection = dbConnection
	dcw.timeKeeper = timeKeeper

	// Eager load the cache
	dcw.updateCache()
	return dcw
}

func (dcw *dynamoCachedWatchDB) getWatchPosition(uuid string) (watchPosition, bool) {
	// If it had been more than 2 hours, always update the cache
	durationSinceUpdate := dcw.timeKeeper.CurrentTime().Sub(dcw.lastUpdatedAt)
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
	dcw.lastUpdatedAt = dcw.timeKeeper.CurrentTime()
}
