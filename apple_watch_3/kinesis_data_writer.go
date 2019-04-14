package apple_watch_3

import (
	"github.com/kine-dmd/api/kinesisqueue"
)

type kinesisDataWriter struct {
	queue kinesisqueue.KinesisQueueInterface
}

func MakeStandardKinesisDataWriter() *kinesisDataWriter {
	// Open a kinesis queue & dynamo DB connection
	const STREAM_NAME = "apple-watch-3"
	queue := kinesisqueue.MakeKinesisQueue(STREAM_NAME)
	return MakeKinesisDataWriter(queue)
}

func MakeKinesisDataWriter(queue kinesisqueue.KinesisQueueInterface) *kinesisDataWriter {
	dataWriter := new(kinesisDataWriter)
	dataWriter.queue = queue
	return dataWriter
}

func (kdw kinesisDataWriter) writeData(data UnparsedAppleWatch3Data) error {
	return kdw.queue.SendToQueue(data, data.WatchPosition.PatientID)
}
