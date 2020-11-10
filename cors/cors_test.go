package cors

import (
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"
	"github.com/stretchr/testify/assert"
)

func Test_keysToString(t *testing.T) {
	tests := []struct {
		name                   string
		input                  map[string]bool
		expectOutputToContain  []string
		expectedNUmberOfCommas int
	}{
		{
			name:                   "empty map returns empty string",
			input:                  nil,
			expectOutputToContain:  nil,
			expectedNUmberOfCommas: 0,
		},
		{
			name:                   "add key from map to string in output",
			input:                  map[string]bool{"hello": true},
			expectOutputToContain:  []string{"hello"},
			expectedNUmberOfCommas: 0,
		},
		{
			name: "All true map is returned in string",
			input: map[string]bool{
				"bob":     true,
				"the":     true,
				"builder": true,
			},
			expectOutputToContain:  []string{"bob", "the", "builder"},
			expectedNUmberOfCommas: 2,
		},
		{
			name: "Only return keys with value true",
			input: map[string]bool{
				"bob":     true,
				"the":     false,
				"builder": true,
			},
			expectOutputToContain:  []string{"bob", "builder"},
			expectedNUmberOfCommas: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commaSeperateTrueKeys(tt.input)

			for _, want := range tt.expectOutputToContain {
				assert.Contains(t, got, want)
			}
			numberOfCommas := strings.Count(got, ", ")
			assert.Equal(t, tt.expectedNUmberOfCommas, numberOfCommas, "actual output was: %q", got)
		})
	}
}

func TestDomains_getAllowedMethodsForOrigin(t *testing.T) {
	tests := []struct {
		name           string
		allowedDomains map[string]map[string]bool
		origin         string
		want           map[string]bool
	}{
		{
			name: "Origin not in list of domains",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			origin: "goodbye",
			want:   nil,
		},
		{
			name: "No methods allowed for origin, return nil",
			allowedDomains: map[string]map[string]bool{
				"hello": nil,
			},
			origin: "hello",
			want:   nil,
		},
		{
			name: "Return methods for origin when there are some",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true, "DELETE": false},
			},
			origin: "hello",
			want:   map[string]bool{"GET": true, "DELETE": false},
		},
		{
			name: "still works with https:// prefix",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true, "DELETE": false},
			},
			origin: "https://hello",
			want:   map[string]bool{"GET": true, "DELETE": false},
		},
		{
			name: "still works with http:// prefix",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true, "DELETE": false},
			},
			origin: "http://hello",
			want:   map[string]bool{"GET": true, "DELETE": false},
		},
		{
			name: "still works with trailing slash",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true, "DELETE": false},
			},
			origin: "hello/",
			want:   map[string]bool{"GET": true, "DELETE": false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Domains{
				Allowed: tt.allowedDomains,
			}

			got := d.getAllowedMethodsForOrigin(tt.origin)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDomains_IsOriginAllowedToPerformMethod(t *testing.T) {
	tests := []struct {
		name           string
		allowedDomains map[string]map[string]bool
		origin         string
		method         string
		want           bool
	}{
		{
			name: "Origin is allowed",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			origin: "hello",
			method: "GET",
			want:   true,
		},
		{
			name: "method not in map",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			origin: "hello",
			method: "DELETE",
			want:   false,
		},
		{
			name: "method not allowed",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": false},
			},
			origin: "hello",
			method: "GET",
			want:   false,
		},
		{
			name: "no methods for origin allowed",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": false},
			},
			origin: "goodbye",
			method: "GET",
			want:   false,
		},
		{
			name: "works with https:// prefix",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			origin: "https://hello",
			method: "GET",
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Domains{
				Allowed: tt.allowedDomains,
			}

			got := d.isOriginAllowedToPerformMethod(tt.origin, tt.method)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDomains_Options(t *testing.T) {
	tests := []struct {
		name           string
		allowedDomains map[string]map[string]bool
		headers        map[string]string
		method         string
		want           events.APIGatewayV2HTTPResponse
	}{
		{
			name: "Origin header not set",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			headers: nil,
			want: events.APIGatewayV2HTTPResponse{
				StatusCode: 400,
				Body:       `{"error": "Origin header was not set on request"}`,
			},
		},
		{
			name: "Return all allowed methods for origin",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			headers: map[string]string{"Origin": "hello"},
			want: events.APIGatewayV2HTTPResponse{
				StatusCode: 204,
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "hello",
					"Access-Control-Allow-Methods": "GET",
				},
			},
		},
		{
			name: "Don't return not allowed methods for origin",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true, "DELETE": false},
			},
			headers: map[string]string{"Origin": "hello"},
			want: events.APIGatewayV2HTTPResponse{
				StatusCode: 204,
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "hello",
					"Access-Control-Allow-Methods": "GET",
				},
			},
		},
		{
			name: "If origin isn't allowed return empty map",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true, "DELETE": false},
			},
			headers: map[string]string{"Origin": "goodbye"},
			want: events.APIGatewayV2HTTPResponse{
				StatusCode: 204,
			},
		},
		{
			name: "Return all allowed methods for origin with https:// prefix",
			allowedDomains: map[string]map[string]bool{
				"hello": {"GET": true},
			},
			headers: map[string]string{"Origin": "https://hello"},
			want: events.APIGatewayV2HTTPResponse{
				StatusCode: 204,
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "https://hello",
					"Access-Control-Allow-Methods": "GET",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Domains{
				Allowed: tt.allowedDomains,
			}

			request := testhelpers.CreateAPIGatewayV2HTTPRequest("does-not-matter", tt.method, "")
			request.Headers = tt.headers

			got := d.Options(request)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDomains_GetCorsHeaders(t *testing.T) {
	tests := []struct {
		name            string
		allowedDomains  map[string]map[string]bool
		origin          string
		method          string
		expectedHeaders map[string]string
		expectedAllowed bool
	}{
		{
			name:            "origin header not set, allowed",
			expectedAllowed: true,
		},
		{
			name: "origin is allowed",
			allowedDomains: map[string]map[string]bool{
				"our-origin": {"GET": true},
			},
			origin: "our-origin",
			method: "GET",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "our-origin",
			},
			expectedAllowed: true,
		},
		{
			name: "origin still allowed with https:// prefix",
			allowedDomains: map[string]map[string]bool{
				"our-origin": {"GET": true},
			},
			origin: "https://our-origin",
			method: "GET",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "https://our-origin",
			},
			expectedAllowed: true,
		},
		{
			name: "not allowed origin returns false",
			allowedDomains: map[string]map[string]bool{
				"our-origin": {"GET": true},
			},
			origin:          "not-ours-m8",
			method:          "GET",
			expectedAllowed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Domains{
				Allowed: tt.allowedDomains,
			}

			request := testhelpers.CreateAPIGatewayV2HTTPRequest("some-path", tt.method, "")
			if tt.origin != "" {
				request.Headers = make(map[string]string)
				request.Headers["Origin"] = tt.origin
			}

			gotHeaders, gotAllowed := d.GetCorsHeaders(request)

			assert.Equal(t, tt.expectedHeaders, gotHeaders)
			assert.Equal(t, tt.expectedAllowed, gotAllowed)
		})
	}
}
