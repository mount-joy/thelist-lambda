package cors

import (
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const accessControlMaxAge = "600" //10 minutes

type OriginChecker interface {
	Options(events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse
	GetCorsHeaders(request events.APIGatewayV2HTTPRequest) map[string]string
}

// Domains contains the allowed domains and the accepted methods for each one
type Domains struct {
	Allowed map[string]map[string]bool
}

// NewOriginChecker returns the default Domains
func NewOriginChecker() OriginChecker {
	acceptedDomains := map[string]map[string]bool{
		"thelist.app":    {http.MethodDelete: true, http.MethodGet: true, http.MethodPatch: true, http.MethodPost: true},
		"dev.thelist.app": {http.MethodDelete: true, http.MethodGet: true, http.MethodPatch: true, http.MethodPost: true},
		"localhost:3000": {http.MethodDelete: true, http.MethodGet: true, http.MethodPatch: true, http.MethodPost: true},
	}
	return &Domains{Allowed: acceptedDomains}
}

func IsOptionsRequest(request events.APIGatewayV2HTTPRequest) bool {
	return http.MethodOptions == request.RequestContext.HTTP.Method
}

// Options performs an options request
// response details methods which are permitted to be performed from the domain in the Origin header
func (d *Domains) Options(request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	origin, ok := caseIncensitiveLookup(originHeader, request.Headers)
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

	headers := headersForOrigin(origin)
	headers[allowMethodsHeader] = commaSeperateTrueKeys(allowedMethds)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNoContent,
		Headers:    headers,
	}
}

func (d *Domains) GetCorsHeaders(request events.APIGatewayV2HTTPRequest) map[string]string {
	origin, ok := caseIncensitiveLookup(originHeader, request.Headers)
	if !ok {
		return nil
	}

	if !d.isOriginAllowedToPerformMethod(origin, request.RequestContext.HTTP.Method) {
		return nil
	}

	return headersForOrigin(origin)
}

func headersForOrigin(origin string) map[string]string {
	return map[string]string{
		allowOriginHeader:  origin,
		maxAgeHeader:       accessControlMaxAge,
		allowHeadersHeader: "content-type",
	}
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

func caseIncensitiveLookup(lookUp string, headers map[string]string) (string, bool) {
	h := http.Header{}
	for key, value := range headers {
		h.Add(key, value)
	}

	value := h.Get(lookUp)
	if value == "" {
		return "", false
	}

	return value, true
}
