package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type dynamoDB struct {
	session dynamodbiface.DynamoDBAPI
}

func createInstance() DB {
	config := aws.Config{Endpoint: aws.String(dbEndpoint)}
	session, err := session.NewSession(&config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create dynamodb session: %s", err.Error()))
	}
	return &dynamoDB{session: dynamodb.New(session)}
}

var instance DB = createInstance()

// DynamoDB returns a databse session using dynamodb
func DynamoDB() DB {
	return instance
}
