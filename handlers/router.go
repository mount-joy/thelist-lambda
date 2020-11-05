package handlers

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/handlers/getitems"
	"github.com/mount-joy/thelist-lambda/handlers/helloworld"
	"github.com/mount-joy/thelist-lambda/handlers/iface"
	"github.com/mount-joy/thelist-lambda/handlers/patchitem"
)

type router struct {
	routes []iface.RouteHandler
}

// NewRouter return the default implementation of Router
func NewRouter() iface.Router {
	routes := []iface.RouteHandler{
		getitems.New(),
		helloworld.New(),
		patchitem.New(),
	}
	return &router{routes: routes}
}

// Route call the appropriate handler for a request based on its path
func (r *router) Route(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	for _, route := range r.routes {
		if route.Match(request) {
			return route.Handle(request)
		}
	}

	log.Printf("Unable to match %s %s", request.RequestContext.HTTP.Method, request.RequestContext.HTTP.Path)
	return nil, http.StatusNotFound
}
