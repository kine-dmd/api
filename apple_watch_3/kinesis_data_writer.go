package apple_watch_3

import (
	"github.com/kine-dmd/api/kinesisqueue"
)

type kinesisDataWriter struct {
	queue kinesisqueue.KinesisQueueInterface
}

func makeStandardKinesisDataWriter() *kinesisDataWriter {
	// Open a kinesis queue & dynamo DB connection
	const STREAM_NAME = "apple-watch-3"
	queue := kinesisqueue.MakeKinesisQueue(STREAM_NAME)
	return makeKinesisDataWriter(queue)
}

func makeKinesisDataWriter(queue kinesisqueue.KinesisQueueInterface) *kinesisDataWriter {
	dataWriter := new(kinesisDataWriter)
	dataWriter.queue = queue
	return dataWriter
}

func (kdw kinesisDataWriter) writeData(data UnparsedAppleWatch3Data) error {
	// Kinesis cannot handle larger than 1MB. Deduct 30% to account for JSON encoding.
	const sizeLimit int = ROW_SIZE_BYTES * 8000

	// Kinesis can only handle 1MB items
	if len(data.RawData) > sizeLimit {

		// Split the first MB of the data off
		firstHalf := UnparsedAppleWatch3Data{data.WatchPosition, data.RawData[:sizeLimit]}

		// Send the first MB of data. If error then return
		err := kdw.queue.SendToQueue(firstHalf, firstHalf.WatchPosition.PatientID)
		if err != nil {
			return err
		}

		// May need to split again so recursively call
		data.RawData = data.RawData[sizeLimit:]
		return kdw.writeData(data)
	}

	// Send the remaining data and return the error
	return kdw.queue.SendToQueue(data, data.WatchPosition.PatientID)
}
