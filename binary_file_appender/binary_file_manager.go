package binary_file_appender

import (
	"sync"
)

type BinaryFileManager interface {
	AppendToFile(filename string, data []byte) error
}

type LocalBinaryFileManager struct {
	fileLocks  map[string]*sync.Mutex
	mapLock    sync.RWMutex
	fileWriter BinaryFileAppender
}

func MakeStandardBinaryFileManager() *LocalBinaryFileManager {
	// Make a file manager using the standard write to file system
	return MakeBinaryFileManager(new(osFileAppender))
}

func MakeBinaryFileManager(appender BinaryFileAppender) *LocalBinaryFileManager {
	// Make a file manager using the provided appender
	manager := new(LocalBinaryFileManager)
	manager.fileWriter = appender

	// Initialise locks
	manager.mapLock = sync.RWMutex{}
	manager.fileLocks = make(map[string]*sync.Mutex)
	return manager
}

func (bfm *LocalBinaryFileManager) AppendToFile(filename string, data []byte) error {
	// Get a pointer to the relevant file lock
	fileLock := bfm.getRelevantFileLock(filename)
	fileLock.Lock()
	defer fileLock.Unlock()

	// Write the data to the file and return any errors
	return bfm.fileWriter.appendToFile(filename, data)
}

func (bfm *LocalBinaryFileManager) getRelevantFileLock(filename string) *sync.Mutex {
	// Acquire the relevant locks - add them if necessary
	bfm.mapLock.RLock()
	fLock, exists := bfm.fileLocks[filename]
	bfm.mapLock.RUnlock()

	// If lock is not already there, then acquire write locks and add it. Need to release read lock to acquire write lock
	if !exists {
		return bfm.addLockForFile(filename)
	}
	return fLock
}

func (bfm *LocalBinaryFileManager) addLockForFile(filename string) *sync.Mutex {
	// Add a new lock for a new file. Need to be sure no-one else is trying to add same lock at same time
	bfm.mapLock.Lock()
	defer bfm.mapLock.Unlock()
	bfm.fileLocks[filename] = new(sync.Mutex)
	return bfm.fileLocks[filename]
}
