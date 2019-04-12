package watch_position_db

import (
	"github.com/kine-dmd/api/api_time"
	"sync"
	"time"
)

// An in memory cache of a UUID -> watch position database
type dynamoCachedWatchDB struct {
	dbConnection  WatchPositionDatabase
	cache         map[string]WatchPosition
	cacheMutex    sync.RWMutex
	lastUpdatedAt time.Time
	timeMutex     sync.RWMutex
	timeKeeper    api_time.ApiTime
}

func MakeStandardDynamoCachedWatchDB() *dynamoCachedWatchDB {
	// Use a dynamo database and the standard system time
	return MakeCachedWatchDB(makeStandardDynamoWatchDatabase(), api_time.SystemTime{})
}

func MakeCachedWatchDB(dbConnection WatchPositionDatabase, timeKeeper api_time.ApiTime) *dynamoCachedWatchDB {
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

func (dcw *dynamoCachedWatchDB) GetTableScan() map[string]WatchPosition {
	// If it had been more than 2 hours, always update the cache
	if dcw.shouldUpdateCache() {
		dcw.updateCache()
	}

	// Obtain a read lock for the cache
	defer dcw.cacheMutex.RUnlock()
	dcw.cacheMutex.RLock()

	// Return value and whether it exists or not
	return dcw.cache
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

	// Get the write lock for the timestamp
	defer dcw.timeMutex.Unlock()
	dcw.timeMutex.Lock()

	// Update the cached values
	dcw.cache = dcw.dbConnection.GetTableScan()
	dcw.lastUpdatedAt = dcw.timeKeeper.CurrentTime()
}
