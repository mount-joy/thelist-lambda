package db

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mount-joy/thelist-lambda/config"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestGetItemsOnList(t *testing.T) {
	testConfig := config.Config{
		Endpoint: "db://thelist",
		TableNames: config.TableNames{
			Items: "items-table",
			Lists: "lists-table",
		},
	}
	tests := []struct {
		name        string
		output      *dynamodb.QueryOutput
		outputErr   error
		expectedRes *[]data.Item
		expectedErr error
	}{
		{
			name:        "When the list is empty no items are returned",
			output:      &dynamodb.QueryOutput{Items: []map[string]*dynamodb.AttributeValue{}},
			outputErr:   nil,
			expectedRes: &[]data.Item{},
			expectedErr: nil,
		},
		{
			name: "When the list contains items, they are returned",
			output: &dynamodb.QueryOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"ListId": {S: aws.String("474c2Fff7")},
						"Item":   {S: aws.String("Oranges")},
						"Id":     {S: aws.String("1c2fa0a1")},
					},
					{
						"ListId": {S: aws.String("474c2Fff7")},
						"Item":   {S: aws.String("Apples")},
						"Id":     {S: aws.String("bb0d5e8e")},
					},
				},
			},
			outputErr: nil,
			expectedRes: &[]data.Item{
				{
					Item:   "Oranges",
					ID:     "1c2fa0a1",
					ListID: "474c2Fff7",
				},
				{
					Item:   "Apples",
					ID:     "bb0d5e8e",
					ListID: "474c2Fff7",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "When Query returns an error, that error is returned",
			output:      &dynamodb.QueryOutput{},
			outputErr:   errors.New("Something went wrong"),
			expectedRes: nil,
			expectedErr: errors.New("Something went wrong"),
		},
		{
			name:        "When Query returns an nil, that nil is returned",
			output:      nil,
			outputErr:   nil,
			expectedRes: nil,
			expectedErr: errors.New("Failed to fetch items"),
		},
	}
	listID := "474c2Fff7"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := dynamodb.QueryInput{
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":id": {S: &listID}},
				KeyConditionExpression:    aws.String("ListId = :id"),
				TableName:                 aws.String("items-table"),
			}
			dbMocked.
				On("Query", &input).
				Return(tt.output, tt.outputErr).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig}

			gotRes, gotErr := d.GetItemsOnList(&listID)

			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
