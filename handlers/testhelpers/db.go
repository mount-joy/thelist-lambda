package testhelpers

import (
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/mock"
)

// MockDB is a Mock of the Database interface
type MockDB struct {
	mock.Mock
}

// CreateItem mocks the DB CreateItem method
func (m *MockDB) CreateItem(listID string, name string) (*data.Item, error) {
	args := m.Called(listID, name)
	return args.Get(0).(*data.Item), args.Error(1)
}

// CreateList mocks the DB CreateList method
func (m *MockDB) CreateList(listName string) (*data.List, error) {
	args := m.Called(listName)
	return args.Get(0).(*data.List), args.Error(1)
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
func (m *MockDB) UpdateItem(listID string, itemID string, newName string, isCompleted *bool) (*data.Item, error) {
	args := m.Called(listID, itemID, newName, isCompleted)
	return args.Get(0).(*data.Item), args.Error(1)
}
