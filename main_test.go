package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
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
			router: handlers.NewRouter(),
		}

		request := events.APIGatewayProxyRequest{
			Path:                  "/hello",
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

func (mr *mockRouter) Route(request events.APIGatewayProxyRequest) (interface{}, int) {
	args := mr.Called(request)
	return args.Get(0), args.Int(1)
}

func TestHandler(t *testing.T) {
	type mockRoute struct {
		body   interface{}
		status int
	}
	tests := []struct {
		name           string
		request        events.APIGatewayProxyRequest
		mockRoute      mockRoute
		expectedBody   string
		expectedStatus int
		expectedErr    error
	}{
		{
			name: "Router doesn't error",
			request: events.APIGatewayProxyRequest{
				Path: "/test",
			},
			mockRoute: mockRoute{
				body:   map[string]string{"message": fmt.Sprintf("huge success")},
				status: 200,
			},
			expectedBody:   "{\"message\":\"huge success\"}",
			expectedStatus: 200,
		},
		{
			name: "Route returns nil",
			request: events.APIGatewayProxyRequest{
				Path: "/test",
			},
			mockRoute: mockRoute{
				body:   nil,
				status: 203,
			},
			expectedStatus: 203,
			expectedBody:   "null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := &mockRouter{}
			router.Test(t)
			defer router.AssertExpectations(t)
			router.
				On("Route", tt.request).
				Return(tt.mockRoute.body, tt.mockRoute.status).
				Once()

			h := handler{
				router: router,
			}

			gotRes, gotErr := h.doRequest(tt.request)

			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedBody, gotRes.Body)
			assert.Equal(t, tt.expectedStatus, gotRes.StatusCode)
		})
	}
}
