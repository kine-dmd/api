package apple_watch_3

import (
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/kinesisqueue"
	"testing"
)

func TestSendingToQueue(t *testing.T) {
	// Create mocks and fake data
	mockCtrl, mockQueue := makeMockQueue(t)
	watchData := makeFakeUnparsedDataStruct("dmd01", 1, make([]byte, ROW_SIZE_BYTES))
	mockQueue.EXPECT().SendToQueue(watchData, watchData.WatchPosition.PatientID).Return(nil).Times(1)

	// Make a kinesis data writer and send to it
	kinesisDataWriter := MakeKinesisDataWriter(mockQueue)
	_ = kinesisDataWriter.writeData(watchData)

	// Check expectations have been satisfied
	mockCtrl.Finish()
}

func TestSplittingLargeItems(t *testing.T) {
	// Create mocks and fake data
	mockCtrl, mockQueue := makeMockQueue(t)
	watchData := makeFakeUnparsedDataStruct("dmd01", 1, make([]byte, ROW_SIZE_BYTES*10000))
	mockQueue.EXPECT().SendToQueue(gomock.Any(), watchData.WatchPosition.PatientID).Return(nil).Times(2)

	// Make a kinesis data writer and send to it
	kinesisDataWriter := MakeKinesisDataWriter(mockQueue)
	_ = kinesisDataWriter.writeData(watchData)

	// Check expectations have been satisfied
	mockCtrl.Finish()
}

func TestTripleSplit(t *testing.T) {
	// Create mocks and fake data
	mockCtrl, mockQueue := makeMockQueue(t)
	watchData := makeFakeUnparsedDataStruct("dmd01", 1, make([]byte, ROW_SIZE_BYTES*18000))
	mockQueue.EXPECT().SendToQueue(integerRowsMatcher{}, watchData.WatchPosition.PatientID).Return(nil).Times(3)

	// Make a kinesis data writer and send to it
	kinesisDataWriter := MakeKinesisDataWriter(mockQueue)
	_ = kinesisDataWriter.writeData(watchData)

	// Check expectations have been satisfied
	mockCtrl.Finish()
}

func makeMockQueue(t *testing.T) (*gomock.Controller, *kinesisqueue.MockKinesisQueueInterface) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockQueue := kinesisqueue.NewMockKinesisQueueInterface(mockCtrl)
	return mockCtrl, mockQueue
}

// GoMock matcher to check that an unparsedAppleWatch3Data struct has an integer number of rows
type integerRowsMatcher struct{}

func (integerRowsMatcher) Matches(x interface{}) bool {
	// Check the data has an integer number of rows
	unparsedStruct := x.(UnparsedAppleWatch3Data)
	return len(unparsedStruct.RawData)%ROW_SIZE_BYTES == 0
}

func (integerRowsMatcher) String() string {
	return "Checks that an unparsed data struct has an integer number or rows (splitting data has not split rows)."
}
