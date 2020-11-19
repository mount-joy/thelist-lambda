package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/cors"
	"github.com/mount-joy/thelist-lambda/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFullFlow(t *testing.T) {
	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		h := handler{
			router:         handlers.NewRouter(),
			allowedDomains: cors.NewOriginChecker(),
		}

		request := events.APIGatewayV2HTTPRequest{
			Headers: map[string]string{"Origin": "thelist.app"},
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: "GET",
					Path:   "/hello",
				},
			},
			QueryStringParameters: map[string]string{"name": "Joy"},
		}

		gotResponse, gotErr := h.doRequest(request)

		expected := "{\"message\":\"Hello, Joy\"}"
		assert.NoError(t, gotErr)
		assert.Equal(t, expected, gotResponse.Body)
	})
}

type mockRouter struct {
	mock.Mock
}

func (mr *mockRouter) Route(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	args := mr.Called(request)
	return args.Get(0), args.Int(1)
}

type mockOriginChecker struct {
	mock.Mock
}

func (mcd *mockOriginChecker) Options(request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	args := mcd.Called(request)
	return args.Get(0).(events.APIGatewayV2HTTPResponse)
}

func (mcd *mockOriginChecker) GetCorsHeaders(request events.APIGatewayV2HTTPRequest) map[string]string {
	args := mcd.Called(request)
	return args.Get(0).(map[string]string)
}

func TestHandler(t *testing.T) {
	type mockRoute struct {
		body   interface{}
		status int
	}
	type mockOptions struct {
		response events.APIGatewayV2HTTPResponse
	}
	type mockGetCorsHeaders struct {
		headers map[string]string
	}
	tests := []struct {
		name               string
		request            events.APIGatewayV2HTTPRequest
		mockOptions        *mockOptions
		mockGetCorsHeaders *mockGetCorsHeaders
		mockRoute          *mockRoute
		expectedBody       string
		expectedStatus     int
		expectedHeaders    map[string]string
		expectedErr        error
	}{
		{
			name: "Router doesn't error",
			request: events.APIGatewayV2HTTPRequest{
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Path: "/test",
					},
				},
			},
			mockGetCorsHeaders: &mockGetCorsHeaders{},
			mockRoute: &mockRoute{
				body:   map[string]string{"message": "huge success"},
				status: 200,
			},
			expectedBody:   "{\"message\":\"huge success\"}",
			expectedStatus: 200,
		},
		{
			name: "Route returns nil",
			request: events.APIGatewayV2HTTPRequest{
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Path: "/test",
					},
				},
			},
			mockGetCorsHeaders: &mockGetCorsHeaders{},
			mockRoute: &mockRoute{
				body:   nil,
				status: 203,
			},
			expectedStatus: 203,
			expectedBody:   "null",
		},
		{
			name: "Sets Access-Control-Allow-Origin when Origin Header is set",
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{"Origin": "test-place"},
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Path:   "/test",
						Method: "GET",
					},
				},
			},
			mockGetCorsHeaders: &mockGetCorsHeaders{
				headers: map[string]string{
					"Access-Control-Allow-Origin": "test-place",
				},
			},
			mockRoute: &mockRoute{
				body:   map[string]string{"message": "huge success"},
				status: 200,
			},
			expectedBody:   "{\"message\":\"huge success\"}",
			expectedStatus: 200,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "test-place",
			},
		},
		{
			name: "If origin isn't allowed, don't add cors headers",
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{"Origin": "some-not-allowed-domain"},
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Path:   "/test",
						Method: "GET",
					},
				},
			},
			mockGetCorsHeaders: &mockGetCorsHeaders{
				headers: nil,
			},
			mockRoute: &mockRoute{
				body:   map[string]string{"message": "huge success"},
				status: 200,
			},
			expectedBody:    "{\"message\":\"huge success\"}",
			expectedStatus:  200,
			expectedHeaders: nil,
		},
		{
			name: "OPTIONS request for allowed domain returns methods",
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{"Origin": "test-place"},
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Method: "OPTIONS",
					},
				},
			},
			mockOptions: &mockOptions{
				response: events.APIGatewayV2HTTPResponse{
					Headers: map[string]string{
						"Access-Control-Allow-Methods": "DELETE, GET, PATCH, POST",
						"Access-Control-Allow-Origin":  "test-place",
					},
					StatusCode: 204,
				},
			},
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods": "DELETE, GET, PATCH, POST",
				"Access-Control-Allow-Origin":  "test-place",
			},
			expectedStatus: 204,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := &mockRouter{}
			router.Test(t)
			defer router.AssertExpectations(t)
			if tt.mockRoute != nil {
				router.
					On("Route", tt.request).
					Return(tt.mockRoute.body, tt.mockRoute.status).
					Once()
			}

			originChecker := &mockOriginChecker{}
			originChecker.Test(t)
			defer originChecker.AssertExpectations(t)
			if tt.mockOptions != nil {
				originChecker.On("Options", tt.request).
					Return(tt.mockOptions.response).
					Once()
			}
			if tt.mockGetCorsHeaders != nil {
				originChecker.On("GetCorsHeaders", tt.request).
					Return(tt.mockGetCorsHeaders.headers).
					Once()
			}

			h := handler{
				router:         router,
				allowedDomains: originChecker,
			}

			gotRes, gotErr := h.doRequest(tt.request)

			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedBody, gotRes.Body)
			assert.Equal(t, tt.expectedStatus, gotRes.StatusCode)
			assert.Equal(t, tt.expectedHeaders, gotRes.Headers)
		})
	}
}
