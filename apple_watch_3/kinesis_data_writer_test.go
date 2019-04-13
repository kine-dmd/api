package apple_watch_3

import (
	"github.com/golang/mock/gomock"
	"github.com/kine-dmd/api/kinesisqueue"
	"github.com/kine-dmd/api/watch_position_db"
	"testing"
)

func TestSendingToQueue(t *testing.T) {
	// Create mocks and fake data
	mockCtrl, mockQueue := makeMockQueue(t)
	watchData := makeFakeUnparsedDataStruct()
	mockQueue.EXPECT().SendToQueue(watchData, watchData.WatchPosition.PatientID).Return(nil).Times(1)

	// Make a kinesis data writer and send to it
	kinesisDataWriter := MakeKinesisDataWriter(mockQueue)
	_ = kinesisDataWriter.writeData(watchData)

	// Check expectations have been satisfied
	mockCtrl.Finish()
}

func makeFakeUnparsedDataStruct() UnparsedAppleWatch3Data {
	watchData := UnparsedAppleWatch3Data{
		watch_position_db.WatchPosition{
			"dmd01", 1,
		},
		[]byte{1, 2},
	}
	return watchData
}

func makeMockQueue(t *testing.T) (*gomock.Controller, *kinesisqueue.MockKinesisQueueInterface) {
	// Make a mock for the kinesis queue
	mockCtrl := gomock.NewController(t)
	mockQueue := kinesisqueue.NewMockKinesisQueueInterface(mockCtrl)
	return mockCtrl, mockQueue
}
