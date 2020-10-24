package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/mount-joy/thelist-lambda/config"
)

type dynamoDB struct {
	session dynamodbiface.DynamoDBAPI
	conf    config.Config
}

func createInstance() DB {
	conf := config.GetConfiguration()
	config := aws.Config{Endpoint: aws.String(conf.Endpoint)}
	return &dynamoDB{session: dynamodb.New(session.New(&config)), conf: conf}
}

var instance DB = createInstance()

// DynamoDB returns a databse session using dynamodb
func DynamoDB() DB {
	return instance
}
