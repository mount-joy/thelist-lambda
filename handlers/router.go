package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Router runs appropriate handler for given path on the request
func Router(request events.APIGatewayProxyRequest) (interface{}, int) {
	if request.Path == "/hello" {
		return helloWorld(request)
	}
	return nil, http.StatusNotFound
}
