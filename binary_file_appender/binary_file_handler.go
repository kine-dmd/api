package binary_file_appender

import (
	"log"
	"os"
)

type BinaryFileAppender interface {
	appendToFile(filename string, data []byte) error
}

type osFileAppender struct{}

func (*osFileAppender) appendToFile(filename string, data []byte) error {
	// Open the file or create it if it doesn't already exist
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("Error opening files %s : %s", filename, err)
		return err
	}

	// Write the data to the file, check for errors
	_, err = f.Write(data)
	if err != nil {
		log.Printf("Error appending to file %s : %s", filename, err)
		return err
	}

	// Try and close the file and check for errors
	err = f.Close()
	if err != nil {
		log.Printf("Error closing file %s : %s", filename, err)
		return err
	}

	return nil
}
