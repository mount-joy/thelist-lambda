package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) UpdateItem(listID string, itemID string, newName string) (*data.Item, error) {
	item := &data.Item{
		ListID: listID,
		ID:     itemID,
		Name:   newName,
	}
	itemToInsert, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	tableName := d.conf.TableNames.Items
	if len(tableName) == 0 {
		panic("Items table name not set")
	}

	input := &dynamodb.PutItemInput{
		Item:                itemToInsert,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_exists(Id) AND attribute_exists(ListId)"),
	}

	_, err = d.session.PutItem(input)

	switch e := err.(type) {
	case nil:
		break
	case awserr.Error:
		if e.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil, ErrorNotFound
		}
	default:
		return nil, err
	}
	return item, nil
}
