package main

import (
	"encoding/json"
	"net/http"

	"github.com/mount-joy/thelist-lambda/cors"
	"github.com/mount-joy/thelist-lambda/handlers"
	"github.com/mount-joy/thelist-lambda/handlers/iface"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type handler struct {
	router         iface.Router
	allowedDomains cors.OriginChecker
}

func (h *handler) doRequest(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if cors.IsOptionsRequest(request) {
		return h.allowedDomains.Options(request), nil
	}

	responseHeaders, allowed := h.allowedDomains.GetCorsHeaders(request)
	if !allowed {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusForbidden,
		}, nil
	}

	result, statusCode := h.router.Route(request)

	res, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			Body:       `{"error": "` + err.Error() + `"}`,
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(res),
		StatusCode: statusCode,
		Headers:    responseHeaders,
	}, nil
}

func main() {
	h := handler{
		router:         handlers.NewRouter(),
		allowedDomains: cors.NewOriginChecker(),
	}

	lambda.Start(h.doRequest)
}
