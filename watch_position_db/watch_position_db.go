package watch_position_db

import (
	"github.com/kine-dmd/api/dynamoDB"
	"log"
	"time"
)

type WatchPosition struct {
	PatientID string `json:"PatientID"`
	Limb      uint8  `json:"Limb"`
}

type WatchPositionDB interface {
	GetWatchPosition(uuid string) (WatchPosition, bool)
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
	cache         map[string]WatchPosition
	lastUpdatedAt time.Time
	timeKeeper    ApiTime
}

func MakeStandardDynamoCachedWatchDB() *dynamoCachedWatchDB {
	// Connect to Dynamo
	const table_name string = "apple_watch_3_positions"
	dbConnection := &dynamoDB.DynamoDBClient{}
	err := dbConnection.InitConn(table_name)
	if err != nil {
		log.Println("Error establishing connection to DynamoDB")
		log.Fatal(err)
	}

	return MakeDynamoCachedWatchDB(dbConnection, systemTime{})
}

func MakeDynamoCachedWatchDB(dbConnection dynamoDB.DynamoDBInterface, timeKeeper ApiTime) *dynamoCachedWatchDB {
	// Create a new connection
	dcw := new(dynamoCachedWatchDB)
	dcw.dbConnection = dbConnection
	dcw.timeKeeper = timeKeeper

	// Eager load the cache
	dcw.updateCache()
	return dcw
}

func (dcw *dynamoCachedWatchDB) GetWatchPosition(uuid string) (WatchPosition, bool) {
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
	var parsedRows = make(map[string]WatchPosition)
	for _, row := range unparsedRows {
		parsedRows[row["uuid"].(string)] = WatchPosition{
			row["patientId"].(string),
			uint8(row["limb"].(float64)),
		}
	}

	// Update the cached values
	dcw.cache = parsedRows
	dcw.lastUpdatedAt = dcw.timeKeeper.CurrentTime()
}
