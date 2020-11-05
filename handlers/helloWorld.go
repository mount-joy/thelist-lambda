package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type helloWorld struct{}

func newHelloWorld() RouteHandler {
	return &helloWorld{}
}

func (h *helloWorld) match(request events.APIGatewayV2HTTPRequest) bool {
	return request.RequestContext.HTTP.Path == "/hello"
}

func (h *helloWorld) handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	name := request.QueryStringParameters["name"]
	return map[string]string{"message": fmt.Sprintf("Hello, %v", name)}, http.StatusOK
}
