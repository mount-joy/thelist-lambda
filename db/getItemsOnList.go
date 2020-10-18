package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *db) GetItemsOnList(listID *string) (*[]data.Item, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {S: listID},
		},
		KeyConditionExpression: aws.String("ListId = :id"),
		TableName:              aws.String(dbTableNameItems),
	}

	result, err := d.session.Query(input)
	if err != nil {
		return nil, err
	}
	if result.Items == nil {
		return nil, nil
	}

	items := []data.Item{}

	for _, i := range result.Items {
		item := new(data.Item)
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return &items, nil
}
