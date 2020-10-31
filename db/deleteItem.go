package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (d *dynamoDB) DeleteItem(listID *string, itemID *string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ListId": {S: listID},
			"Id":     {S: itemID},
		},
		TableName: aws.String(d.conf.TableNames.Items),
	}

	_, err := d.session.DeleteItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return NewError(ErrorNotFound)
			}
		}
	}
	return err
}
