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
	tests := []struct {
		name            string
		mockResponseErr error
		expectedErr     error
	}{
		{
			name:            "If the item exists it is deleted",
			mockResponseErr: nil,
			expectedErr:     nil,
		},
		{
			name:            "When db returns an error, that error is returned",
			mockResponseErr: errors.New("Something went wrong"),
			expectedErr:     errors.New("Something went wrong"),
		},
		{
			name:            "When Query returns condition not match error, then succeeds",
			mockResponseErr: awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedErr:     nil,
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
				TableName: stringToPointer("items-table"),
			}
			dbMocked.
				On("DeleteItem", &input).
				Return(&dynamodb.DeleteItemOutput{}, tt.mockResponseErr).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig}
			gotErr := d.DeleteItem(listID, itemID)

			assert.Equal(t, tt.expectedErr, gotErr)
		})
	}
}
