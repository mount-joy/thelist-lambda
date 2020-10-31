package iface

import "github.com/aws/aws-lambda-go/events"

// Router - interface for routing requests to the right handler
type Router interface {
	Route(events.APIGatewayV2HTTPRequest) (interface{}, int)
}

// RouteHandler - interface for matching and handling a particular request
type RouteHandler interface {
	Match(events.APIGatewayV2HTTPRequest) bool
	Handle(events.APIGatewayV2HTTPRequest) (interface{}, int)
}
