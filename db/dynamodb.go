package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type dynamoDB struct {
	session dynamodbiface.DynamoDBAPI
}

// NewDynamoDB returns a new database session using dynamodb
func NewDynamoDB() DB {
	config := aws.Config{Endpoint: aws.String(dbEndpoint)}
	return &dynamoDB{session: dynamodb.New(session.New(&config))}
}
