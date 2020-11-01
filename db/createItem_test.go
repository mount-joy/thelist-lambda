package db

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	listID := "474c2Fff7"
	itemID := "b6cf642d"
	tableName := "items-table"
	conditionExpression := "attribute_not_exists(Id)"
	itemName := "Peaches"
	tests := []struct {
		name           string
		outputErr      error
		expectedOutput *data.Item
		expectedErr    error
	}{
		{
			name:           "If the ID does not exists it creates the item",
			outputErr:      nil,
			expectedOutput: &data.Item{ID: itemID, ListID: listID, Name: itemName},
			expectedErr:    nil,
		},
		{
			name:        "When db returns an error, that error is returned",
			outputErr:   errors.New("Something went wrong"),
			expectedErr: errors.New("Something went wrong"),
		},
		{
			name:        "When DB returns condition not match error, not found error is returned",
			outputErr:   awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedErr: ErrorIDExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := dynamodb.PutItemInput{
				Item: map[string]*dynamodb.AttributeValue{
					"Id":     {S: &itemID},
					"ListId": {S: &listID},
					"Name":   {S: &itemName},
				},
				TableName:           &tableName,
				ConditionExpression: &conditionExpression,
			}
			dbMocked.
				On("PutItem", &input).
				Return(&dynamodb.PutItemOutput{}, tt.outputErr).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig, generateID: func() string { return itemID }}
			gotRes, gotErr := d.CreateItem(&listID, &itemName)

			assert.Equal(t, tt.expectedOutput, gotRes)
			assert.Equal(t, tt.expectedErr, gotErr)
		})
	}
}
