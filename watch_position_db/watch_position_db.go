package watch_position_db

import (
	"github.com/kine-dmd/api/api_time"
	"github.com/kine-dmd/api/dynamoDB"
	"log"
	"sync"
	"time"
)

type WatchPosition struct {
	PatientID string `json:"PatientID"`
	Limb      uint8  `json:"Limb"`
}

type WatchPositionDB interface {
	GetWatchPosition(uuid string) (WatchPosition, bool)
}

type dynamoCachedWatchDB struct {
	dbConnection  dynamoDB.DynamoDBInterface
	cache         map[string]WatchPosition
	cacheMutex    sync.RWMutex
	lastUpdatedAt time.Time
	timeMutex     sync.RWMutex
	timeKeeper    api_time.ApiTime
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

	return MakeDynamoCachedWatchDB(dbConnection, api_time.SystemTime{})
}

func MakeDynamoCachedWatchDB(dbConnection dynamoDB.DynamoDBInterface, timeKeeper api_time.ApiTime) *dynamoCachedWatchDB {
	// Create a new connection
	dcw := new(dynamoCachedWatchDB)
	dcw.dbConnection = dbConnection
	dcw.timeKeeper = timeKeeper

	// Eager load the cache
	dcw.cacheMutex = sync.RWMutex{}
	dcw.timeMutex = sync.RWMutex{}
	dcw.updateCache()
	return dcw
}

func (dcw *dynamoCachedWatchDB) GetWatchPosition(uuid string) (WatchPosition, bool) {
	// If it had been more than 2 hours, always update the cache
	if dcw.shouldUpdateCache() {
		dcw.updateCache()
	}

	// Obtain a read lock for the cache
	defer dcw.cacheMutex.RUnlock()
	dcw.cacheMutex.RLock()

	// Return value and whether it exists or not
	val, ok := dcw.cache[uuid]
	return val, ok
}

func (dcw *dynamoCachedWatchDB) shouldUpdateCache() bool {
	// Get a read lock for the time
	defer dcw.timeMutex.RUnlock()
	dcw.timeMutex.RLock()

	// If it had been more than 2 hours update the cache
	durationSinceUpdate := dcw.timeKeeper.CurrentTime().Sub(dcw.lastUpdatedAt)
	return durationSinceUpdate.Hours() >= 2
}

func (dcw *dynamoCachedWatchDB) updateCache() {
	// Get the timestamp for the cache
	defer dcw.cacheMutex.Unlock()
	dcw.cacheMutex.Lock()

	// Check someone else has not already updated the cache while we waited
	if !dcw.shouldUpdateCache() {
		return
	}

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

	// Get the write lock for the timestamp
	defer dcw.timeMutex.Unlock()
	dcw.timeMutex.Lock()

	// Update the cached values
	dcw.cache = parsedRows
	dcw.lastUpdatedAt = dcw.timeKeeper.CurrentTime()
}
