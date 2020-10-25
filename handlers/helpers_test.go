package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) GetItemsOnList(input *string) (*[]data.Item, error) {
	args := m.Called(input)
	return args.Get(0).(*[]data.Item), args.Error(1)
}

func (m *mockDB) UpdateItem(listID *string, itemID *string, newName *string) (*data.Item, error) {
	args := m.Called(listID, itemID, newName)
	return args.Get(0).(*data.Item), args.Error(1)
}

func createAPIGatewayV2HTTPRequest(path string, method string, body string) events.APIGatewayV2HTTPRequest {
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
