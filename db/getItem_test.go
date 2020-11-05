package db

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestGetItem(t *testing.T) {
	listID := "474c2Fff7"
	itemID := "b6cf642d"
	name := "Cheese"
	tableName := "items-table"
	tests := []struct {
		name          string
		mockOutputErr error
		mockOutput    *dynamodb.GetItemOutput
		expectedRes   *data.Item
		expectedErr   error
	}{
		{
			name:          "If the item exists it is retrieved",
			mockOutputErr: nil,
			mockOutput: &dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{
					"Id":     {S: &itemID},
					"ListId": {S: &listID},
					"Name":   {S: &name},
				},
			},
			expectedRes: &data.Item{ID: itemID, ListID: listID, Name: name},
			expectedErr: nil,
		},
		{
			name:          "When db returns an error, that error is returned",
			mockOutputErr: errors.New("Something went wrong"),
			mockOutput:    nil,
			expectedRes:   nil,
			expectedErr:   errors.New("Something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := dynamodb.GetItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Id":     {S: &itemID},
					"ListId": {S: &listID},
				},
				TableName: &tableName,
			}
			dbMocked.
				On("GetItem", &input).
				Return(tt.mockOutput, tt.mockOutputErr).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig}
			gotRes, gotErr := d.GetItem(listID, itemID)

			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
