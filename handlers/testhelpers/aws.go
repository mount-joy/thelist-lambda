package testhelpers

import "github.com/aws/aws-lambda-go/events"

// CreateAPIGatewayV2HTTPRequest is a helper function for creating a APIGatewayV2HTTPRequest object
func CreateAPIGatewayV2HTTPRequest(path string, method string, body string) events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest{
		Body: body,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path:   path,
				Method: method,
			},
		},
	}
}
