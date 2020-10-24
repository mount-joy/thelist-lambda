package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func helloWorld(request events.APIGatewayProxyRequest) (interface{}, int) {
	name := request.QueryStringParameters["name"]
	return map[string]string{"message": fmt.Sprintf("Hello, %v", name)}, http.StatusOK
}
