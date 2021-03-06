package handlers

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/handlers/iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRoute struct {
	mock.Mock
}

func (m *mockRoute) Match(request events.APIGatewayV2HTTPRequest) bool {
	args := m.Called(request)
	return args.Bool(0)
}

func (m *mockRoute) Handle(request events.APIGatewayV2HTTPRequest) (interface{}, int) {
	args := m.Called(request)
	return args.Get(0), args.Int(1)
}

func TestRoute(t *testing.T) {
	bodyA := map[string]string{"route": "A"}
	bodyB := map[string]string{"route": "B"}

	tests := []struct {
		name           string
		matchResA      bool
		matchResB      bool
		expectedBody   interface{}
		expectedStatus int
	}{
		{
			name:           "No routes match then returns 404",
			matchResA:      false,
			matchResB:      false,
			expectedBody:   nil,
			expectedStatus: 404,
		},
		{
			name:           "Route A matches then returns body A",
			matchResA:      true,
			matchResB:      false,
			expectedBody:   bodyA,
			expectedStatus: 200,
		},
		{
			name:           "Route B matches then returns body B",
			matchResA:      false,
			matchResB:      true,
			expectedBody:   bodyB,
			expectedStatus: 200,
		},
		{
			name:           "Both routes matches then returns body A",
			matchResA:      true,
			matchResB:      true,
			expectedBody:   bodyA,
			expectedStatus: 200,
		},
	}

	request := events.APIGatewayV2HTTPRequest{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routeA := &mockRoute{}
			routeA.Test(t)
			routeA.
				On("Match", request).
				Return(tt.matchResA)

			routeA.
				On("Handle", request).
				Return(bodyA, 200)

			routeB := &mockRoute{}
			routeB.Test(t)
			routeB.
				On("Match", request).
				Return(tt.matchResB)

			routeB.
				On("Handle", request).
				Return(bodyB, 200)

			r := router{
				routes: []iface.RouteHandler{routeA, routeB},
			}

			gotRes, gotStatusCode := r.Route(request)

			assert.Equal(t, tt.expectedBody, gotRes)
			assert.Equal(t, tt.expectedStatus, gotStatusCode)
		})
	}
}
