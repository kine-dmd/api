package dynamoDB

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type DynamoDBInterface interface {
	InitConn(tableName string) error
	GetTableScan() []map[string]interface{}
}

type DynamoDBClient struct {
	connection *dynamodb.DynamoDB
	tableName  string
}

// initConn opens the connection to the dynamo DB database
func (db *DynamoDBClient) InitConn(tableName string) error {
	// Save the table name
	db.tableName = tableName

	// Create a session in a given AWS region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	if err != nil {
		log.Println("Got error creating session:")
		log.Println(err.Error())
		return err
	}

	// Create DynamoDB client
	db.connection = dynamodb.New(sess)
	return nil
}

// Get a scan of the entire table
func (db *DynamoDBClient) GetTableScan() []map[string]interface{} {
	// Take a scan of the table
	params := &dynamodb.ScanInput{
		TableName: aws.String(db.tableName),
	}
	result, err := db.connection.Scan(params)

	// Check for errors
	if err != nil {
		log.Println("Got error doing scan:", err.Error())
		return nil
	}

	// Create list to store result in
	var allRows = make([]map[string]interface{}, len(result.Items))

	// Unmarshall to list of maps
	for index, row := range result.Items {
		err = dynamodbattribute.UnmarshalMap(row, &allRows[index])
		if err != nil {
			log.Println("Got error unmarshalling:", err.Error())
			return nil
		}
	}
	return allRows
}
