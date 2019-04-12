package kinesisqueue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"log"
)

type KinesisQueueInterface interface {
	SendToQueue(data interface{}, shardId string) error
}

type KinesisQueueClient struct {
	kinesis    *kinesis.Kinesis
	streamName string
}

// MakeKinesisQueue opens the connection to the location event kinesis queue
func MakeKinesisQueue(streamName string) *KinesisQueueClient {
	// Define the stream name and the AWS region it's in
	region := "eu-west-2"
	// Create a new AWS session in the required region
	s, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatal("Unable to make AWS connection for Kinesis", err.Error())
	}

	// Create a new kinesis adapter (assume stream exists on AWS)
	kq := new(KinesisQueueClient)
	kq.kinesis = kinesis.New(s)
	kq.streamName = streamName

	return kq
}

// Pre: the event object is valid
func (kq *KinesisQueueClient) SendToQueue(data interface{}, shardId string) error {
	// Encode a record into JSON bytes
	byteEncodedData, _ := json.Marshal(data)

	// Send the record to Kinesis
	_, err := kq.kinesis.PutRecord(&kinesis.PutRecordInput{
		Data:         byteEncodedData,
		StreamName:   aws.String(kq.streamName),
		PartitionKey: aws.String(shardId),
	})
	if err != nil {
		log.Println("Error sending item to Kinesis")
		log.Println(err)
		return err
	}
	return nil
}
