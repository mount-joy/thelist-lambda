package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) UpdateItem(listID string, itemID string, newName string, isCompleted *bool) (*data.Item, error) {
	key, err := dynamodbattribute.MarshalMap(&data.ItemKey{ID: itemID, ListID: listID})
	if err != nil {
		return nil, err
	}

	tableName := d.conf.TableNames.Items
	if len(tableName) == 0 {
		panic("Items table name not set")
	}

	fieldsToUpdate, updateExpression, expressionAttributeNames := getUpdateFields(newName, isCompleted)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: fieldsToUpdate,
		Key:                       key,
		TableName:                 aws.String(tableName),
		UpdateExpression:          updateExpression,
		ReturnValues:              aws.String("ALL_NEW"),
		ExpressionAttributeNames:  expressionAttributeNames,
		ConditionExpression:       aws.String("attribute_exists(Id)"),
	}

	output, err := d.session.UpdateItem(input)

	switch e := err.(type) {
	case nil:
		break
	case awserr.Error:
		if e.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil, ErrorNotFound
		}
		if e.Code() == "ValidationException" { // https://github.com/aws/aws-sdk-go/issues/3140
			return nil, ErrorBadRequest
		}
		return nil, err
	default:
		return nil, err
	}

	item := new(data.Item)
	err = dynamodbattribute.UnmarshalMap(output.Attributes, &item)
	return item, err
}

func getUpdateFields(newName string, isCompleted *bool) (map[string]*dynamodb.AttributeValue, *string, map[string]*string) {
	fields := map[string]*dynamodb.AttributeValue{}
	names := map[string]*string{}
	var updateExpression *string

	if isCompleted != nil {
		fields[":c"] = &dynamodb.AttributeValue{BOOL: isCompleted}
		updateExpression = appendUpdateExpression(updateExpression, "IsCompleted = :c")
	}

	if newName != "" {
		fields[":n"] = &dynamodb.AttributeValue{S: aws.String(newName)}
		names["#n"] = aws.String("Name")
		updateExpression = appendUpdateExpression(updateExpression, "#n = :n")
	}

	// If names isn't used it needs to be nil
	if len(names) == 0 {
		names = nil
	}

	return fields, updateExpression, names
}

func appendUpdateExpression(updateExpression *string, newPart string) *string {
	if updateExpression == nil {
		return aws.String(fmt.Sprintf("SET %s", newPart))
	}
	return aws.String(fmt.Sprintf("%s, %s", *updateExpression, newPart))
}
