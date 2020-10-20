package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type dynamoDB struct {
	session dynamodbiface.DynamoDBAPI
}

var instance DB

func createInstance() DB {
	log.Println("Creating new instance!")
	config := aws.Config{Endpoint: aws.String(dbEndpoint)}
	return &dynamoDB{session: dynamodb.New(session.New(&config))}
}

// DynamoDB returns a databse session using dynamodb
func DynamoDB() DB {
	if instance == nil {
		instance = createInstance()
	}
	return instance
}
