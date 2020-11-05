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
		ListID: listID,
		ID:     itemID,
		Name:   name,
	}
	itemToInsert, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:                itemToInsert,
		TableName:           aws.String(d.conf.TableNames.Items),
		ConditionExpression: aws.String("attribute_not_exists(Id)"),
	}

	_, err = d.session.PutItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return nil, ErrorIDExists
			}
		}
		return nil, err
	}
	return item, nil
}
