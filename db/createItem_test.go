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
	itemName := "Peaches"
	timestamp := "2020-01-23T09:59:14.9396531Z"

	tests := []struct {
		name           string
		item           map[string]*dynamodb.AttributeValue
		mockOutputErr  error
		expectedOutput *data.Item
		expectedErr    error
	}{
		{
			name:           "If the ID does not exists it creates the item",
			item:           createExpectedInput(itemID, listID, itemName, false, timestamp),
			mockOutputErr:  nil,
			expectedOutput: &data.Item{ItemKey: data.ItemKey{ID: itemID, ListID: listID}, Name: itemName, IsCompleted: false, UpdatedTimestamp: timestamp, CreatedTimestamp: timestamp},
			expectedErr:    nil,
		},
		{
			name:          "When db returns an error, that error is returned",
			item:          createExpectedInput(itemID, listID, itemName, false, timestamp),
			mockOutputErr: errors.New("Something went wrong"),
			expectedErr:   errors.New("Something went wrong"),
		},
		{
			name:          "When DB returns condition not match error, not found error is returned",
			item:          createExpectedInput(itemID, listID, itemName, false, timestamp),
			mockOutputErr: awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedErr:   ErrorIDExists,
		},
		{
			name:          "When DB unrecognised awserr, passon the error",
			item:          createExpectedInput(itemID, listID, itemName, false, timestamp),
			mockOutputErr: awserr.New("uh oh", "whoops", errors.New("Oh dear")),
			expectedErr:   awserr.New("uh oh", "whoops", errors.New("Oh dear")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := dynamodb.PutItemInput{
				Item:                tt.item,
				TableName:           stringToPointer("items-table"),
				ConditionExpression: stringToPointer("attribute_not_exists(Id)"),
			}
			dbMocked.
				On("PutItem", &input).
				Return(&dynamodb.PutItemOutput{}, tt.mockOutputErr).
				Once()

			d := dynamoDB{
				session:      dbMocked,
				conf:         testConfig,
				generateID:   func() string { return itemID },
				getTimestamp: func() string { return timestamp },
			}
			gotRes, gotErr := d.CreateItem(listID, itemName)

			assert.Equal(t, tt.expectedOutput, gotRes)
			assert.Equal(t, tt.expectedErr, gotErr)
		})
	}
}

func createExpectedInput(itemID string, listID string, itemName string, isCompleted bool, timestamp string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"Id":          {S: &itemID},
		"ListId":      {S: &listID},
		"Name":        {S: &itemName},
		"IsCompleted": {BOOL: &isCompleted},
		"Created":     {S: &timestamp},
		"Updated":     {S: &timestamp},
	}
}
