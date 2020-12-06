package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) CreateItem(listID string, name string) (*data.Item, error) {
	itemID := d.generateID()

	item := &data.Item{
		ItemKey: data.ItemKey{
			ListID: listID,
			ID:     itemID,
		},
		Name:        name,
		IsCompleted: false,
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
		ConditionExpression: aws.String("attribute_not_exists(Id)"),
	}

	_, err = d.session.PutItem(input)

	switch e := err.(type) {
	case nil:
		break
	case awserr.Error:
		if e.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil, ErrorIDExists
		}
		return nil, err
	default:
		return nil, err
	}

	return item, nil
}
