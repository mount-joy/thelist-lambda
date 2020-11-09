package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/mount-joy/thelist-lambda/config"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
	dynamodbiface.DynamoDBAPI
}

func (m *mockDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockDB) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func (m *mockDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *mockDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

var testConfig config.Config = config.Config{
	Endpoint: "db://thelist",
	TableNames: config.TableNames{
		Items: "items-table",
		Lists: "lists-table",
	},
}

func stringToPointer(input string) *string {
	return &input
}
