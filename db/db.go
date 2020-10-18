package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/mount-joy/thelist-lambda/data"
)

// DB - interface for talking to the database
type DB interface {
	GetItemsOnList(listID *string) (*[]data.Item, error)
}

type db struct {
	session dynamodbiface.DynamoDBAPI
}

// NewDB returns a new database connection
func NewDB() DB {
	config := aws.Config{Endpoint: aws.String(dbEndpoint)}
	return &db{session: dynamodb.New(session.New(&config))}
}
