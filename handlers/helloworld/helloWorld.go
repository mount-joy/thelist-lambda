package helloworld

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/handlers/iface"
)

type helloWorld struct{}

// New returns an instance of helloWorld satisfying the RouteHandler interface
func New() iface.RouteHandler {
	return &helloWorld{}
}

func (h *helloWorld) Match(request events.APIGatewayV2HTTPRequest) bool {
	return request.RequestContext.HTTP.Path == "/hello"
}

func (h *helloWorld) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	name := request.QueryStringParameters["name"]
	return map[string]string{"message": fmt.Sprintf("Hello, %v", name)}, http.StatusOK
}
