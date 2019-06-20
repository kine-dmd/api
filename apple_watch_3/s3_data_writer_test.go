package apple_watch_3

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/s3Connection"
	"testing"
)

func TestS3DataWriterCallsUpload(t *testing.T) {
	// Create mocks and fake data
	mockCtrl, mockS3 := makeMockS3Connection(t)
	byteData := make([]byte, ROW_SIZE_BYTES*18000)
	watchData := makeFakeUnparsedDataStruct("dmd01", 1, byteData)
	mockS3.EXPECT().UploadFile(BUCKET_NAME, gomock.Any(), bytes.NewReader(byteData))

	// Make a kinesis data writer and send to it
	kinesisDataWriter := makeS3DataWriter(mockS3)
	_ = kinesisDataWriter.writeData(watchData)

	// Check expectations have been satisfied
	mockCtrl.Finish()
}

func makeMockS3Connection(t *testing.T) (*gomock.Controller, *s3Connection.MockS3Uploader) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockS3 := s3Connection.NewMockS3Uploader(mockCtrl)
	return mockCtrl, mockS3
}
