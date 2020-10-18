package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type helloWorld struct{}

// NewRouter return the default implementation of the hello world route handler
func newHelloWorld() RouteHandler {
	return &helloWorld{}
}

func (h *helloWorld) match(request events.APIGatewayProxyRequest) bool {
	return request.Path == "/hello"
}

func (h *helloWorld) handle(request events.APIGatewayProxyRequest) (interface{}, int) {
	name := request.QueryStringParameters["name"]
	return map[string]string{"message": fmt.Sprintf("Hello, %v", name)}, http.StatusOK
}
