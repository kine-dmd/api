package apple_watch_3

import (
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/binary_file_appender"
	"testing"
)

func TestCallsFileManager(t *testing.T) {
	// Make mocks and local file data writer
	mockCtrl, mockFileManager := makeFileManagerMocks(t)
	writer := MakeLocalFileDataWriter(mockFileManager)

	// Make a standard fake data struct
	unparsedDataStruct := makeFakeUnparsedDataStruct()
	mockFileManager.EXPECT().AppendToFile("~/data/dmd01/leftHand.bin", []byte{1, 2}).Return(nil).Times(1)

	// Write the data
	err := writer.writeData(unparsedDataStruct)

	// Check expectations
	if err != nil {
		t.Fatalf("Got unexpected error: %s", err)
	}
	mockCtrl.Finish()
}

func makeFileManagerMocks(t *testing.T) (*gomock.Controller, *binary_file_appender.MockBinaryFileManager) {
	// Make a mock file manager
	mockCtrl := gomock.NewController(t)
	mockFileManager := binary_file_appender.NewMockBinaryFileManager(mockCtrl)
	return mockCtrl, mockFileManager
}
