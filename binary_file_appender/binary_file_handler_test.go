package binary_file_appender

import (
	"os"
	"testing"
)

func TestCreatesFileWhenNoneExist(t *testing.T) {
	// Use a unique filename for each test
	filename := "TestCreatesFileWhenNoneExist.bin"
	defer func() { _ = os.Remove(filename) }()

	// Try creating and writing to the file
	err := new(osFileAppender).appendToFile(filename, []byte{})
	if err != nil {
		t.Fatalf("Unable to append to originally non-existant file: %s", err)
	}

	checkFileProperties(t, filename, 0)
}

func TestAppendingToFile(t *testing.T) {
	// Try appending to a file that does not yet exist
	filename := "TestAppendingToFile.bin"
	defer func() { _ = os.Remove(filename) }()

	fileAppender := new(osFileAppender)
	err := fileAppender.appendToFile(filename, []byte{0, 0, 0})
	if err != nil {
		t.Fatalf("Unable to append to originally non-existant file: %s", err)
	}

	// Try appending to a file that does exist
	err = fileAppender.appendToFile(filename, []byte{1, 1})
	if err != nil {
		t.Fatalf("Unable to append to existing file: %s", err)
	}

	// Delete the file to cleanup tests
	checkFileProperties(t, filename, 5)
}

func checkFileProperties(t *testing.T, filename string, expectedSize int) {
	// Stat the file to check it exists
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("File does not exist after first write: %s", err)
	}

	if info.Size() != int64(expectedSize) {
		t.Fatalf("Filesize does not match. Expected %d. Got %d.", expectedSize, info.Size())
	}
}
