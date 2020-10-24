package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Router - interface for routing requests to the right handler
type Router interface {
	Route(events.APIGatewayProxyRequest) (interface{}, int)
}

// RouteHandler - interface for matching and handling a particular request
type RouteHandler interface {
	match(request events.APIGatewayProxyRequest) bool
	handle(events.APIGatewayProxyRequest) (interface{}, int)
}

type router struct {
	routes []RouteHandler
}

// NewRouter return the default implementation of Router
func NewRouter() Router {
	routes := []RouteHandler{
		newHelloWorld(),
		newGetItems(),
		newPostList(),
	}
	return &router{routes: routes}
}

// Route call the appropriate handler for a request based on its path
func (r *router) Route(request events.APIGatewayProxyRequest) (interface{}, int) {
	for _, route := range r.routes {
		if route.match(request) {
			return route.handle(request)
		}
	}

	return nil, http.StatusNotFound
}
