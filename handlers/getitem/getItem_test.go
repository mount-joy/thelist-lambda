package getitem

import (
	"testing"

	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestGetItemMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "GET",
			expectedRes: true,
		},
		{
			name:        "Returns true for a uppercase ID",
			path:        "/lists/B6CF642D/items/73BB82C4/",
			method:      "GET",
			expectedRes: true,
		},
		{
			name:        "Returns true without trailing slash",
			path:        "/lists/b6cf642d/items/73bb82c4",
			method:      "GET",
			expectedRes: true,
		},
		{
			name:        "Returns false for list path",
			path:        "/lists/b6cf642d/",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for list path without trailing slash",
			path:        "/lists/b6cf642d",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path",
			path:        "/lists/",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path without trailing slash",
			path:        "/lists",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for items path",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6/items",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false when path is empty",
			path:        "",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for a POST request",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false for a DELETE request",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false for a PATCH request",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "PATCH",
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &testhelpers.MockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, tt.method, "")
			d := getItems{db: dbMocked}
			gotRes := d.Match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

type mockGetItem struct {
	res *data.Item
	err error
}

func TestGetItemHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		itemID             string
		mockOutput         *mockGetItem
		expectedRes        interface{}
		expectedStatusCode int
	}{
		{
			name:               "Returns 'Internal Server Error' when the path is empty",
			path:               "",
			listID:             "test-list-id",
			expectedRes:        nil,
			expectedStatusCode: 500,
		},
		{
			name:               "Returns 'Internal Server Error' when the path is not in the correct format",
			path:               "/lists/test-list-id",
			listID:             "test-list-id",
			expectedRes:        nil,
			expectedStatusCode: 500,
		},
		{
			name:               "Returns 'OK' and results when the path matches",
			path:               "/lists/test-list-id/items/test-item-id",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			mockOutput:         &mockGetItem{res: &data.Item{Name: "ABC", ID: "888"}, err: nil},
			expectedRes:        &data.Item{Name: "ABC", ID: "888"},
			expectedStatusCode: 200,
		},
		{
			name:               "Returns 'Internal Server Error' when the path is not in the correct format",
			path:               "/lists/test-list-id/items",
			listID:             "test-list-id",
			expectedRes:        nil,
			expectedStatusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &testhelpers.MockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			if tt.mockOutput != nil {
				dbMocked.
					On("GetItem", tt.listID, tt.itemID).
					Return(tt.mockOutput.res, tt.mockOutput.err).
					Once()
			}

			d := getItems{db: dbMocked}

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, "GET", "")
			gotRes, statusCode := d.Handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
