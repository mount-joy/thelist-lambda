package db

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestCreateList(t *testing.T) {
	listID := "1234"
	tests := []struct {
		name           string
		listName       string
		mockOutputErr  error
		expectedOutput *data.List
		expectedErr    error
	}{
		{
			name:           "If dynamodb passes, creates the list",
			listName:       "my-list",
			mockOutputErr:  nil,
			expectedOutput: &data.List{ListKey: data.ListKey{ID: listID}, Name: "my-list"},
			expectedErr:    nil,
		},
		{
			name:           "If dynamodb failes, pass back the error",
			listName:       "my-list",
			mockOutputErr:  fmt.Errorf("not working"),
			expectedOutput: nil,
			expectedErr:    fmt.Errorf("not working"),
		},
		{
			name:           "If there is a clash in dynamodb, return ErrorIDExists",
			listName:       "my-list",
			mockOutputErr:  awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedOutput: nil,
			expectedErr:    ErrorIDExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			item := map[string]*dynamodb.AttributeValue{
				"Id":   {S: &listID},
				"Name": {S: &tt.listName},
			}
			input := dynamodb.PutItemInput{
				Item:                item,
				TableName:           stringToPointer("lists-table"),
				ConditionExpression: stringToPointer("attribute_not_exists(Id)"),
			}
			dbMocked.
				On("PutItem", &input).
				Return(&dynamodb.PutItemOutput{}, tt.mockOutputErr).
				Once()

			d := dynamoDB{
				session:    dbMocked,
				conf:       testConfig,
				generateID: func() string { return listID },
			}

			gotRes, gotErr := d.CreateList(tt.listName)

			assert.Equal(t, tt.expectedOutput, gotRes)
			assert.Equal(t, tt.expectedErr, gotErr)
		})
	}
}
