package kinesisqueue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"log"
	"os"
)

type KinesisQueueInterface interface {
	InitConn(streamName string) error
	SendToQueue(data interface{}, shardId string) error
}

type KinesisQueueClient struct {
	kinesis    *kinesis.Kinesis
	streamName string
}

// InitConn opens the connection to the location event kinesis queue
func (kq *KinesisQueueClient) InitConn(streamName string) error {
	// Define the stream name and the AWS region it's in
	region := "eu-west-2"
	// Create a new AWS session in the required region
	s, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	// Create a new kinesis adapter (assume stream exists
	kq.kinesis = kinesis.New(s)
	kq.streamName = streamName

	return nil
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
