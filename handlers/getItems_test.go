package handlers

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) GetItemsOnList(input *string) (*[]data.Item, error) {
	args := m.Called(input)
	return args.Get(0).(*[]data.Item), args.Error(1)
}

func TestGetItemsMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6/items/",
			expectedRes: true,
		},
		{
			name:        "Returns true for a uppercase ID",
			path:        "/lists/4BBA7AB4-1D3A-4694-990B-6F78DEFC84C1/items/",
			expectedRes: true,
		},
		{
			name:        "Returns true without trailing slash",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6/items",
			expectedRes: true,
		},
		{
			name:        "Returns false for list path",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6/",
			expectedRes: false,
		},
		{
			name:        "Returns false for list path without trailing slash",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path",
			path:        "/lists/",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path without trailing slash",
			path:        "/lists",
			expectedRes: false,
		},
		{
			name:        "Returns false for item path",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6/items/14b58681-b5c3-4e5a-a928-12e9dc63cdb3",
			expectedRes: false,
		},
		{
			name:        "Returns false when path is empty",
			path:        "",
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := events.APIGatewayProxyRequest{Path: tt.path}
			d := getItems{db: dbMocked}
			gotRes := d.match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

func TestGetItemsHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		output             *[]data.Item
		outputErr          error
		expectedRes        interface{}
		expectedStatusCode int
		shouldCallDB       bool
	}{
		{
			name:               "Returns 'Internal Server Error' when the path is empty",
			path:               "",
			listID:             "test-list-id",
			output:             nil,
			outputErr:          nil,
			expectedRes:        nil,
			expectedStatusCode: 500,
			shouldCallDB:       false,
		},
		{
			name:               "Returns 'Internal Server Error' when the path is not in the correct format",
			path:               "/lists/test-list-id",
			listID:             "test-list-id",
			output:             nil,
			outputErr:          nil,
			expectedRes:        nil,
			expectedStatusCode: 500,
			shouldCallDB:       false,
		},
		{
			name:               "Returns 'OK' and results when the path matches",
			path:               "/lists/test-list-id/items",
			listID:             "test-list-id",
			output:             &[]data.Item{data.Item{Item: "ABC", ID: "888"}},
			outputErr:          nil,
			expectedRes:        &[]data.Item{data.Item{Item: "ABC", ID: "888"}},
			expectedStatusCode: 200,
			shouldCallDB:       true,
		},
		{
			name:               "Returns 'Internal Server Error' when the path is not in the correct format",
			path:               "/lists/test-list-id/items",
			listID:             "test-list-id",
			output:             nil,
			outputErr:          errors.New("It went wrong"),
			expectedRes:        nil,
			expectedStatusCode: 500,
			shouldCallDB:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			if tt.shouldCallDB {
				defer dbMocked.AssertExpectations(t)
			}

			dbMocked.
				On("GetItemsOnList", &tt.listID).
				Return(tt.output, tt.outputErr)

			d := getItems{db: dbMocked}

			input := events.APIGatewayProxyRequest{
				Path: tt.path,
			}
			gotRes, statusCode := d.handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
