package main

import (
	"encoding/json"

	"github.com/mount-joy/thelist-lambda/handlers"
	"github.com/mount-joy/thelist-lambda/handlers/iface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type handler struct {
	router iface.Router
}

func (h *handler) doRequest(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	result, statusCode := h.router.Route(request)

	res, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body:       "{\"error\": \"" + err.Error() + "\"}",
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
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
