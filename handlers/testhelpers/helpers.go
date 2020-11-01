package testhelpers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/mock"
)

// MockDB is a Mock of the Database interface
type MockDB struct {
	mock.Mock
}

// CreateItem mocks the DB CreateItem method
func (m *MockDB) CreateItem(listID *string, name *string) (*data.Item, error) {
	args := m.Called(listID, name)
	return args.Get(0).(*data.Item), args.Error(1)
}

// GetItem mocks the DB GetItem method
func (m *MockDB) GetItem(listID string, itemID string) (*data.Item, error) {
	args := m.Called(listID, itemID)
	return args.Get(0).(*data.Item), args.Error(1)
}

// DeleteItem mocks the DB DeleteItem method
func (m *MockDB) DeleteItem(listID string, itemID string) error {
	args := m.Called(listID, itemID)
	return args.Error(1)
}

// GetItemsOnList mocks the DB GetItemsOnList method
func (m *MockDB) GetItemsOnList(input string) (*[]data.Item, error) {
	args := m.Called(input)
	return args.Get(0).(*[]data.Item), args.Error(1)
}

// UpdateItem mocks the DB UpdateItem method
func (m *MockDB) UpdateItem(listID string, itemID string, newName string) (*data.Item, error) {
	args := m.Called(listID, itemID, newName)
	return args.Get(0).(*data.Item), args.Error(1)
}

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
