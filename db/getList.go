package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) GetList(listID string) (*data.List, error) {
	tableName := d.conf.TableNames.Lists
	if len(tableName) == 0 {
		panic("Items table name not set")
	}

	key, err := dynamodbattribute.MarshalMap(data.ListKey{ID: listID})

	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}
	res, err := d.session.GetItem(input)

	if err != nil {
		return nil, err
	}

	item := new(data.List)
	err = dynamodbattribute.UnmarshalMap(res.Item, &item)
	return item, err
}
