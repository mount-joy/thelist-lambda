package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Router - interface for routing requests to the right handler
type Router interface {
	Route(events.APIGatewayProxyRequest) (interface{}, int)
}

type router struct{}

// NewRouter return the default implementation of Router
func NewRouter() Router {
	return &router{}
}

// Route call the appropriate handler for a request based on its path
func (r *router) Route(request events.APIGatewayProxyRequest) (interface{}, int) {
	if request.Path == "/hello" {
		return helloWorld(request)
	}
	return nil, http.StatusNotFound
}
