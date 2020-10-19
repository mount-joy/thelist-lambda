package main

import (
	"encoding/json"

	"github.com/mount-joy/thelist-lambda/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type handler struct {
	router handlers.Router
}

func (h *handler) doRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, statusCode := h.router.Route(request)

	res, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: statusCode,
	}, nil
}

func main() {
	h := handler{
		router: handlers.NewRouter(),
	}

	lambda.Start(h.doRequest)
}
