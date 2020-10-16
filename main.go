package main

import (
	"encoding/json"

	"github.com/mount-joy/thelist-lambda/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, statusCode := handlers.Router(request)
	res, err := json.Marshal(result)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: statusCode,
	}, nil
}

func main() {
	lambda.Start(handler)
}
