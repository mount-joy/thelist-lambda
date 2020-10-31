package db

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestDeleteItem(t *testing.T) {
	listID := "474c2Fff7"
	itemID := "b6cf642d"
	tableName := "items-table"
	tests := []struct {
		name        string
		outputErr   error
		expectedErr error
	}{
		{
			name:        "If the item exists it is deleted",
			outputErr:   nil,
			expectedErr: nil,
		},
		{
			name:        "When db returns an error, that error is returned",
			outputErr:   errors.New("Something went wrong"),
			expectedErr: errors.New("Something went wrong"),
		},
		{
			name:        "When Query returns condition not match error, not found error is returned",
			outputErr:   awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedErr: NewError(ErrorNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := dynamodb.DeleteItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Id":     {S: &itemID},
					"ListId": {S: &listID},
				},
				TableName: &tableName,
			}
			dbMocked.
				On("DeleteItem", &input).
				Return(&dynamodb.DeleteItemOutput{}, tt.outputErr).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig}
			gotErr := d.DeleteItem(&listID, &itemID)

			assert.Equal(t, tt.expectedErr, gotErr)
		})
	}
}
