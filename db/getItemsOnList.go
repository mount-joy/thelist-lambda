package db

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) GetItemsOnList(listID string) (*[]data.Item, error) {
	tableName := d.conf.TableNames.Items
	if len(tableName) == 0 {
		panic("Items table name not set")
	}

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {S: &listID},
		},
		KeyConditionExpression: aws.String("ListId = :id"),
		TableName:              aws.String(tableName),
	}

	result, err := d.session.Query(input)
	if err != nil {
		return nil, err
	}
	if result == nil || result.Items == nil {
		return nil, errors.New("Failed to fetch items")
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
