package s3Connection

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"math/rand"
	"strconv"
	"testing"
)

const (
	TEST_BUCKET_NAME = "kine-dmd-test"
)

func TestS3Upload100(t *testing.T) {
	testS3Upload(t, 100)
}

func TestS3Upload10000(t *testing.T) {
	testS3Upload(t, 10000)
}

func TestS3Upload1000000(t *testing.T) {
	testS3Upload(t, 1000000)
}

func testS3Upload(t *testing.T, fileSizeBytes int) {
	// Make the connection, generate a random file name and random data
	s3Conn := MakeS3Connection()
	filename := strconv.FormatUint(rand.Uint64(), 10)
	originalData := makeRandomData(fileSizeBytes)

	// Upload the data
	err := s3Conn.UploadFile(TEST_BUCKET_NAME, filename, bytes.NewReader(originalData))
	if err != nil {
		t.Fatalf("Error uploading file to S3: %s", err)
	}

	// Download the file and compare it
	retrievedData := retrieveTestFile(t, filename)
	compareData(t, originalData, retrievedData)
}

func retrieveTestFile(t *testing.T, filePath string) []byte {
	// Make a session to download the file
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	))
	downloader := s3manager.NewDownloader(sess)

	// Create a buffer in memory to store the binary data
	buffer := aws.NewWriteAtBuffer([]byte{})

	// Download the file from S3 to the buffer
	_, err := downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(TEST_BUCKET_NAME),
			Key:    aws.String(filePath),
		})
	if err != nil {
		t.Fatalf("Unable to download test file from bucket: %s", err)
	}

	return buffer.Bytes()
}

func compareData(t *testing.T, expected []byte, actual []byte) {
	if len(expected) != len(actual) {
		t.Fatalf("Expected result length not equal to actual. Expected length %d. Got length %d.", len(expected), len(actual))
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("Byte %d does not match. Expected %b, got %b", i, expected[i], actual[i])
		}
	}
}

func makeRandomData(length int) []byte {
	data := make([]byte, length)
	rand.Read(data)
	return data
}
