package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mount-joy/thelist-lambda/data"
)

func (d *dynamoDB) CreateList(listName string) (*data.List, error) {
	timestamp := d.getTimestamp()

	list := &data.List{
		ListKey: data.ListKey{
			ID: d.generateID(),
		},
		Name:             listName,
		CreatedTimestamp: timestamp,
		UpdatedTimestamp: timestamp,
		IsShared:         false,
	}

	listToInsert, err := dynamodbattribute.MarshalMap(list)
	if err != nil {
		return nil, err
	}

	tableName := d.conf.TableNames.Lists
	if len(tableName) == 0 {
		panic("Lists table name not set")
	}
	input := &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                listToInsert,
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

	return list, nil
}
