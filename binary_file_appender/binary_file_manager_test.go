package binary_file_appender

import (
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"
)

func TestManagerWriteCallsFileWrite(t *testing.T) {
	// Make the mocks and the managers
	mockCtrl, mockFileAppender := makeGoMockFileAppender(t)
	manager := MakeBinaryFileManager(mockFileAppender)

	// Query the manager and set mock expectations
	filename := "Somefile"
	mockFileAppender.EXPECT().appendToFile(filename, []byte{}).Return(nil).Times(1)
	err := manager.AppendToFile(filename, []byte{})

	// Check no errors were thrown and expectations were satisfied
	if err != nil {
		t.Fatalf("Write to file returned error %s when it shouldn't have.", err)
	}
	mockCtrl.Finish()
}

func TestManagerSimultaneousWriteCallsFileWrite(t *testing.T) {
	// Make the mocks and the managers
	mockCtrl, mockFileAppender := makeGoMockFileAppender(t)
	manager := MakeBinaryFileManager(mockFileAppender)

	// Set mock expectations
	mockFileAppender.EXPECT().appendToFile("1", []byte{}).Return(nil).Times(1)
	mockFileAppender.EXPECT().appendToFile("2", []byte{}).Return(nil).Times(1)

	// Use two threads to query the manager
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() { _ = manager.AppendToFile("1", []byte{}); wg.Done() }()
	go func() { _ = manager.AppendToFile("2", []byte{}); wg.Done() }()

	// Check no errors were thrown and expectations were satisfied
	wg.Wait()
	mockCtrl.Finish()
}

func TestSimultaneousFileWritesToSingleFile(t *testing.T) {
	// Create a mock file app
	mockFileAppender := simultaneousFileAppendChecker{sync.Map{}, t}
	manager := MakeBinaryFileManager(mockFileAppender)

	// Spawn 100 threads to query the file appender which checks they are being called sequentially
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			_ = manager.AppendToFile("AnyFile", []byte{})
		}()
	}
	wg.Wait()
}

func TestSimultaneousFileWritesToMultipleFiles(t *testing.T) {
	// Create a mock file app
	mockFileAppender := simultaneousFileAppendChecker{sync.Map{}, t}
	manager := MakeBinaryFileManager(mockFileAppender)

	// Spawn 100 threads to query the file appender which checks they are being called sequentially
	wg := sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			_ = manager.AppendToFile("AnyFile", []byte{})
		}()
		go func() {
			defer wg.Done()
			_ = manager.AppendToFile("AnyFile2", []byte{})
		}()
	}
	wg.Wait()
}

func TestSimulataneousWritesTakeSameTime(t *testing.T) {
	// Check how long each takes to run
	startTime := time.Now()
	TestSimultaneousFileWritesToSingleFile(t)
	middleTime := time.Now()
	TestSimultaneousFileWritesToMultipleFiles(t)
	endTime := time.Now()

	// Time taken to run should be roughly equal due to same sleep time and parallel writes to separate files
	timeDifference := (startTime.Sub(middleTime)) - (middleTime.Sub(endTime))
	if timeDifference > time.Millisecond*10 {
		t.Fatalf("Writing to two files using mock files should take "+
			"same amount of time as one. Took %f seconds longer", timeDifference.Seconds())
	}
}

/***********************************************************
 A mock file appender that can be used to check if the same
 file is being written to by two threads simultaneously
***********************************************************/
type simultaneousFileAppendChecker struct {
	activeMap sync.Map
	t         *testing.T
}

func (sfac simultaneousFileAppendChecker) appendToFile(filename string, data []byte) error {
	// Check that another thread has not already made this active
	isActive, exists := sfac.activeMap.Load(filename)
	if exists && isActive.(bool) {
		sfac.t.Fatal("Another thread tried to activate an already active file.")
	}

	// We declare it active
	sfac.activeMap.Store(filename, true)

	// Sleep for 100 milliseconds to simulate writing to file
	time.Sleep(10 * time.Millisecond)

	// No longer using it
	sfac.activeMap.Store(filename, false)
	return nil
}

func makeGoMockFileAppender(t *testing.T) (*gomock.Controller, *MockBinaryFileAppender) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockFileAppender := NewMockBinaryFileAppender(mockCtrl)
	return mockCtrl, mockFileAppender
}
