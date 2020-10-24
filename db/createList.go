package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (d *dynamoDB) CreateList(listID *string) (error) {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": { S: listID }, "Name": { S: aws.String("Shopping List") },
		},
		TableName:              aws.String(dbTableNameItems),
	}

	_, err := d.session.PutItem(input)
	return err
}
