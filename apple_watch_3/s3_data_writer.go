package apple_watch_3

import (
	"bytes"
	"github.com/kine-dmd/api/s3Connection"
	"log"
	"strconv"
	"time"
)

const (
	BUCKET_NAME = "kine-dmd-aw3-intermediary"
)

type s3DataWriter struct {
	s3Conn *s3Connection.S3Connection
}

func MakeStandardS3DataWriter() *s3DataWriter {
	// Open a kinesis queue & dynamo DB connection
	s3dw := new(s3DataWriter)
	s3dw.s3Conn = s3Connection.MakeS3Connection()
	return s3dw
}

func (s3dw s3DataWriter) writeData(data UnparsedAppleWatch3Data) error {
	// Create a new filepath for this file
	filePath := data.WatchPosition.PatientID + "/" + strconv.Itoa(int(data.WatchPosition.Limb)) + "/" + strconv.Itoa(int(time.Now().UnixNano())) + ".bin"

	// Try and upload the data
	err := s3dw.s3Conn.UploadFile(BUCKET_NAME, filePath, bytes.NewReader(data.RawData))
	if err != nil {
		log.Printf("Error writing file %s to S3. %s", filePath, err)
	}

	return err
}
