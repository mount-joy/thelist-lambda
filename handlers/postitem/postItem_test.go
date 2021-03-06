package postitem

import (
	"fmt"
	"testing"

	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestPostItemMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists/b6cf642d/items/",
			method:      "POST",
			expectedRes: true,
		},
		{
			name:        "Returns true for a uppercase ID",
			path:        "/lists/B6CF642D/items/",
			method:      "POST",
			expectedRes: true,
		},
		{
			name:        "Returns true without trailing slash",
			path:        "/lists/b6cf642d/items",
			method:      "POST",
			expectedRes: true,
		},
		{
			name:        "Returns false for list path",
			path:        "/lists/b6cf642d/",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false for list path without trailing slash",
			path:        "/lists/b6cf642d",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path",
			path:        "/lists/",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path without trailing slash",
			path:        "/lists",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false for item path",
			path:        "/lists/b6cf642d-7a72-4969-bcc9-73bb82c4b3f6/items/73bb82c4",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false when path is empty",
			path:        "",
			method:      "POST",
			expectedRes: false,
		},
		{
			name:        "Returns false for a GET request",
			path:        "/lists/b6cf642d/items/",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for a DELETE request",
			path:        "/lists/b6cf642d/items/",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false for a PATCH request",
			path:        "/lists/b6cf642d/items/",
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
			d := postItem{db: dbMocked}
			gotRes := d.Match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

type mockPostItem struct {
	res *data.Item
	err error
}

func TestPostItemHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		itemName           string
		body               string
		mockOutput         *mockPostItem
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
			path:               "/lists/test-list-id/items/",
			listID:             "test-list-id",
			itemName:           "my item",
			body:               "{ \"Name\": \"my item\" }",
			mockOutput:         &mockPostItem{res: &data.Item{Name: "ABC", ItemKey: data.ItemKey{ID: "888"}}, err: nil},
			expectedRes:        &data.Item{Name: "ABC", ItemKey: data.ItemKey{ID: "888"}},
			expectedStatusCode: 200,
		},
		{
			name:     "Returns 'internal server error' if database errors",
			path:     "/lists/test-list-id/items/",
			listID:   "test-list-id",
			itemName: "my item",
			body:     "{ \"Name\": \"my item\" }",
			mockOutput: &mockPostItem{
				res: &data.Item{Name: "ABC", ItemKey: data.ItemKey{ID: "888"}},
				err: fmt.Errorf("broken"),
			},
			expectedRes:        nil,
			expectedStatusCode: 500,
		},
		{
			name:               "Returns 'Bad Request' when the body is empty",
			path:               "/lists/test-list-id/items",
			listID:             "test-list-id",
			body:               "",
			expectedRes:        nil,
			expectedStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &testhelpers.MockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			if tt.mockOutput != nil {
				dbMocked.
					On("CreateItem", tt.listID, tt.itemName).
					Return(tt.mockOutput.res, tt.mockOutput.err).
					Once()
			}

			d := postItem{db: dbMocked}

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, "POST", tt.body)
			gotRes, statusCode := d.Handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
