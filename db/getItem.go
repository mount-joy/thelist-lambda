package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) GetItem(listID string, itemID string) (*data.Item, error) {
	tableName := d.conf.TableNames.Items
	if len(tableName) == 0 {
		panic("Items table name not set")
	}

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ListId": {S: &listID},
			"Id":     {S: &itemID},
		},
		TableName: aws.String(tableName),
	}

	res, err := d.session.GetItem(input)

	if err != nil {
		return nil, err
	}

	item := new(data.Item)
	err = dynamodbattribute.UnmarshalMap(res.Item, &item)
	return item, err
}
