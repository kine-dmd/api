package apple_watch_3

import (
	"github.com/kine-dmd/api/kinesisqueue"
)

type Kinesis_data_writer struct {
	queue kinesisqueue.KinesisQueueInterface
}

func MakeStandardKinesisDataWriter() *Kinesis_data_writer {
	// Open a kinesis queue & dynamo DB connection
	const STREAM_NAME = "apple-watch-3"
	queue := kinesisqueue.MakeKinesisQueue(STREAM_NAME)
	return MakeKinesisDataWriter(queue)
}

func MakeKinesisDataWriter(queue kinesisqueue.KinesisQueueInterface) *Kinesis_data_writer {
	dataWriter := new(Kinesis_data_writer)
	dataWriter.queue = queue
	return dataWriter
}

func (kdw Kinesis_data_writer) writeData(data UnparsedAppleWatch3Data) error {
	return kdw.queue.SendToQueue(data, data.WatchPosition.PatientID)
}
