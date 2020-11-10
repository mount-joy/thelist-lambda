package cors

import (
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type OriginChecker interface {
	Options(events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse
	GetCorsHeaders(request events.APIGatewayV2HTTPRequest) (map[string]string, bool)
}

// Domains contains the allowed domains and the accepted methods for each one
type Domains struct {
	Allowed map[string]map[string]bool
}

// NewDomains returns the default Domains
func NewOriginChecker() OriginChecker {
	acceptedDomains := map[string]map[string]bool{
		"thelist.app": {"DELETE": true, "GET": true, "PATCH": true, "POST": true},
	}
	return &Domains{Allowed: acceptedDomains}
}

func IsOptionsRequest(request events.APIGatewayV2HTTPRequest) bool {
	return http.MethodOptions == request.RequestContext.HTTP.Method
}

// Options performs an options request
// response details methods which are permitted to be performed from the domain in the Origin header
func (d *Domains) Options(request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	origin, ok := request.Headers["Origin"]
	if !ok {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Origin header was not set on request"}`,
		}
	}

	allowedMethds := d.getAllowedMethodsForOrigin(origin)
	if len(allowedMethds) == 0 {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNoContent,
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  origin,
			"Access-Control-Allow-Methods": commaSeperateTrueKeys(allowedMethds),
		},
	}
}

func (d *Domains) GetCorsHeaders(request events.APIGatewayV2HTTPRequest) (map[string]string, bool) {
	origin, ok := request.Headers["Origin"]
	if !ok {
		return nil, true
	}

	if !d.isOriginAllowedToPerformMethod(origin, request.RequestContext.HTTP.Method) {
		return nil, false
	}

	return map[string]string{
		"Access-Control-Allow-Origin": origin,
	}, true
}

func (d *Domains) isOriginAllowedToPerformMethod(origin string, method string) bool {
	allowedMethods := d.getAllowedMethodsForOrigin(origin)
	if len(allowedMethods) == 0 {
		return false
	}

	allowed, ok := allowedMethods[method]
	if !ok {
		log.Printf("domain %q tried use %q method", origin, method)
		return false
	}

	return allowed
}

func (d *Domains) getAllowedMethodsForOrigin(origin string) map[string]bool {
	trimmedOrigin := removePrefixAndSuffix(origin)

	methods, ok := d.Allowed[trimmedOrigin]
	if !ok {
		log.Printf("%q tried to make request", trimmedOrigin)
		return nil
	}

	if len(methods) == 0 {
		log.Printf("No allowed methods for %q", trimmedOrigin)
		return nil
	}

	return methods
}

func commaSeperateTrueKeys(input map[string]bool) string {
	output := ""
	for key, value := range input {
		if !value {
			continue
		}

		output += key + ", "
	}

	return strings.TrimSuffix(output, ", ")
}

func removePrefixAndSuffix(origin string) string {
	output := strings.TrimPrefix(origin, "https://")
	output = strings.TrimPrefix(output, "http://")

	output = strings.TrimSuffix(output, "/")

	return output
}
