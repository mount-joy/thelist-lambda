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
	timestamp := "2020-01-23T09:59:14.9396531Z"

	tests := []struct {
		testName                         string
		newName                          string
		isCompleted                      *bool
		mockedErrResponse                error
		mockedResponse                   *dynamodb.UpdateItemOutput
		expectedUpdateExpression         *string
		expectedFieldsToUpdate           map[string]*dynamodb.AttributeValue
		expectedKey                      map[string]*dynamodb.AttributeValue
		expectedExpressionAttributeNames map[string]*string
		expectedRes                      *data.Item
		expectedErr                      error
	}{
		{
			testName:                         "If the item exists it is updated",
			newName:                          newName,
			isCompleted:                      boolToPointer(true),
			mockedResponse:                   updateItemOutput(listID, itemID, newName, true),
			mockedErrResponse:                nil,
			expectedUpdateExpression:         stringToPointer("SET IsCompleted = :c, #n = :n, Updated = :t"),
			expectedFieldsToUpdate:           updateBothFields(newName, true, timestamp),
			expectedKey:                      map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedExpressionAttributeNames: map[string]*string{"#n": stringToPointer("Name")},
			expectedRes:                      &data.Item{ItemKey: data.ItemKey{ID: itemID, ListID: listID}, Name: newName, IsCompleted: true},
			expectedErr:                      nil,
		},
		{
			testName:                         "If only a new name is supplied, only it is updated",
			newName:                          newName,
			isCompleted:                      nil,
			mockedResponse:                   updateItemOutput(listID, itemID, newName, false),
			mockedErrResponse:                nil,
			expectedUpdateExpression:         stringToPointer("SET #n = :n, Updated = :t"),
			expectedFieldsToUpdate:           updateName(newName, timestamp),
			expectedKey:                      map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedExpressionAttributeNames: map[string]*string{"#n": stringToPointer("Name")},
			expectedRes:                      &data.Item{ItemKey: data.ItemKey{ID: itemID, ListID: listID}, Name: newName, IsCompleted: false},
			expectedErr:                      nil,
		},
		{
			testName:                         "If only isCompleted is changed, only it is updated",
			newName:                          "",
			isCompleted:                      boolToPointer(true),
			mockedResponse:                   updateItemOutput(listID, itemID, newName, true),
			mockedErrResponse:                nil,
			expectedUpdateExpression:         stringToPointer("SET IsCompleted = :c, Updated = :t"),
			expectedFieldsToUpdate:           updateIsCompleted(true, timestamp),
			expectedKey:                      map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedExpressionAttributeNames: nil,
			expectedRes:                      &data.Item{ItemKey: data.ItemKey{ID: itemID, ListID: listID}, Name: newName, IsCompleted: true},
			expectedErr:                      nil,
		},
		{
			testName:                         "When db returns an error, that error is returned",
			newName:                          newName,
			isCompleted:                      boolToPointer(true),
			mockedResponse:                   nil,
			mockedErrResponse:                errors.New("Something went wrong"),
			expectedUpdateExpression:         stringToPointer("SET IsCompleted = :c, #n = :n, Updated = :t"),
			expectedFieldsToUpdate:           updateBothFields(newName, true, timestamp),
			expectedKey:                      map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedExpressionAttributeNames: map[string]*string{"#n": stringToPointer("Name")},
			expectedRes:                      nil,
			expectedErr:                      errors.New("Something went wrong"),
		},
		{
			testName:                         "If the item doesn't exist the Query returns condition not match error, not found error is returned",
			newName:                          newName,
			isCompleted:                      boolToPointer(true),
			mockedResponse:                   nil,
			mockedErrResponse:                awserr.New(dynamodb.ErrCodeConditionalCheckFailedException, "Bad", errors.New("Oh dear")),
			expectedUpdateExpression:         stringToPointer("SET IsCompleted = :c, #n = :n, Updated = :t"),
			expectedFieldsToUpdate:           updateBothFields(newName, true, timestamp),
			expectedKey:                      map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedExpressionAttributeNames: map[string]*string{"#n": stringToPointer("Name")},
			expectedRes:                      nil,
			expectedErr:                      ErrorNotFound,
		},
		{
			testName:                 "If the update request is invalid, BadRequest is returned",
			newName:                  "",
			mockedResponse:           nil,
			mockedErrResponse:        awserr.New("ValidationException", "Bad", errors.New("Oh dear")),
			expectedUpdateExpression: stringToPointer("SET Updated = :t"),
			expectedFieldsToUpdate:   map[string]*dynamodb.AttributeValue{":t": &dynamodb.AttributeValue{S: &timestamp}},
			expectedKey:              map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedRes:              nil,
			expectedErr:              ErrorBadRequest,
		},
		{
			testName:                         "If another AWS error is returned that error message is passed on",
			newName:                          newName,
			isCompleted:                      boolToPointer(true),
			mockedResponse:                   nil,
			mockedErrResponse:                awserr.New("Oops", "Bad", errors.New("Oh dear")),
			expectedUpdateExpression:         stringToPointer("SET IsCompleted = :c, #n = :n, Updated = :t"),
			expectedFieldsToUpdate:           updateBothFields(newName, true, timestamp),
			expectedKey:                      map[string]*dynamodb.AttributeValue{"Id": {S: &itemID}, "ListId": {S: &listID}},
			expectedExpressionAttributeNames: map[string]*string{"#n": stringToPointer("Name")},
			expectedRes:                      nil,
			expectedErr:                      awserr.New("Oops", "Bad", errors.New("Oh dear")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := dynamodb.UpdateItemInput{
				ExpressionAttributeValues: tt.expectedFieldsToUpdate,
				Key:                       tt.expectedKey,
				TableName:                 stringToPointer("items-table"),
				UpdateExpression:          tt.expectedUpdateExpression,
				ReturnValues:              stringToPointer("ALL_NEW"),
				ExpressionAttributeNames:  tt.expectedExpressionAttributeNames,
				ConditionExpression:       stringToPointer("attribute_exists(Id)"),
			}
			dbMocked.
				On("UpdateItem", &input).
				Return(tt.mockedResponse, tt.mockedErrResponse).
				Once()

			d := dynamoDB{session: dbMocked, conf: testConfig, getTimestamp: func() string { return timestamp }}
			gotRes, gotErr := d.UpdateItem(listID, itemID, tt.newName, tt.isCompleted)

			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

func updateItemOutput(listID string, itemID string, name string, isCompleted bool) *dynamodb.UpdateItemOutput {
	return &dynamodb.UpdateItemOutput{
		Attributes: map[string]*dynamodb.AttributeValue{
			"Id": &dynamodb.AttributeValue{
				S: stringToPointer(itemID),
			},
			"ListId": &dynamodb.AttributeValue{
				S: stringToPointer(listID),
			},
			"Name": &dynamodb.AttributeValue{
				S: stringToPointer(name),
			},
			"IsCompleted": &dynamodb.AttributeValue{
				BOOL: boolToPointer(isCompleted),
			},
		},
	}
}

func updateBothFields(name string, isCompleted bool, timestamp string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		":c": &dynamodb.AttributeValue{BOOL: &isCompleted},
		":n": &dynamodb.AttributeValue{S: &name},
		":t": &dynamodb.AttributeValue{S: &timestamp},
	}
}

func updateName(name string, timestamp string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		":n": &dynamodb.AttributeValue{S: &name},
		":t": &dynamodb.AttributeValue{S: &timestamp},
	}
}

func updateIsCompleted(isCompleted bool, timestamp string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		":c": &dynamodb.AttributeValue{BOOL: &isCompleted},
		":t": &dynamodb.AttributeValue{S: &timestamp},
	}
}
