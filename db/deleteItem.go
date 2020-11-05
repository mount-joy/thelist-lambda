package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (d *dynamoDB) DeleteItem(listID string, itemID string) error {
	tableName := d.conf.TableNames.Items
	if len(tableName) == 0 {
		panic("Items table name not set")
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ListId": {S: &listID},
			"Id":     {S: &itemID},
		},
		TableName: aws.String(tableName),
	}

	_, err := d.session.DeleteItem(input)

	switch e := err.(type) {
	case nil:
		break
	case awserr.Error:
		if e.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
			return nil
		}
	default:
		return err
	}

	return nil
}
