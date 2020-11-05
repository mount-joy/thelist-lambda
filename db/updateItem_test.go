package db

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestUpdateItem(t *testing.T) {
	listID := "474c2Fff7"
	itemID := "b6cf642d"
	newName := "Cheese"
	tableName := "items-table"
	conditionExpression := "attribute_exists(Id) AND attribute_exists(ListId)"
	tests := []struct {
		name              string
		mockedErrResponse error
		expectedRes       *data.Item
		expectedErr       error
	}{
		{
			name:              "If the item exists it is updated",
			mockedErrResponse: nil,
			expectedRes:       &data.Item{ID: itemID, ListID: listID, Name: newName},
			expectedErr:       nil,
		},
		{
			name:              "When db returns an error, that error is returned",
			mockedErrResponse: errors.New("Something went wrong"),
			expectedRes:       nil,
			expectedErr:       errors.New("Something went wrong"),
		},
		{
			name:              "If the item doesn't exist the Query returns condition not match error, not found error is returned",
			mockedErrResponse: awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedRes:       nil,
			expectedErr:       ErrorNotFound,
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
					"Name":   {S: &newName},
				},
				TableName:           &tableName,
				ConditionExpression: &conditionExpression,
			}
			dbMocked.
				On("PutItem", &input).
				Return(&dynamodb.PutItemOutput{}, tt.mockedErrResponse).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig}
			gotRes, gotErr := d.UpdateItem(listID, itemID, newName)

			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
