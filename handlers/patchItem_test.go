package handlers

import (
	"errors"
	"testing"

	"github.com/mount-joy/thelist-lambda/db"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestPatchItemMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "PATCH",
			expectedRes: true,
		},
		{
			name:        "Returns true for a uppercase ID",
			path:        "/lists/B6CF642D/items/73BB82C4/",
			method:      "PATCH",
			expectedRes: true,
		},
		{
			name:        "Returns true without trailing slash",
			path:        "/lists/b6cf642d/items/73bb82c4",
			method:      "PATCH",
			expectedRes: true,
		},
		{
			name:        "Returns false for list path",
			path:        "/lists/b6cf642d/",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false for list path without trailing slash",
			path:        "/lists/b6cf642d",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path",
			path:        "/lists/",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path without trailing slash",
			path:        "/lists",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false for items path",
			path:        "/lists/b6cf642d/items/",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false when path is empty",
			path:        "",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false for a GET request",
			path:        "/lists/b6cf642d/items/73bb82c4/",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := &mockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := createAPIGatewayV2HTTPRequest(tt.path, tt.method, "")
			d := patchItem{db: dbMocked}
			gotRes := d.match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

func TestPatchItemHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		itemID             string
		newName            string
		body               string
		output             *data.Item
		outputErr          error
		expectedRes        interface{}
		expectedStatusCode int
		shouldCallDB       bool
	}{
		{
			name:               "Returns 'Bad Request' when the path is empty",
			path:               "",
			listID:             "test-list-id",
			output:             nil,
			outputErr:          nil,
			expectedRes:        nil,
			expectedStatusCode: 400,
			shouldCallDB:       false,
		},
		{
			name:               "Returns 'Bad Request' when the path is not in the correct format",
			path:               "/lists/test-list-id",
			listID:             "test-list-id",
			output:             nil,
			outputErr:          nil,
			expectedRes:        nil,
			expectedStatusCode: 400,
			shouldCallDB:       false,
		},
		{
			name:               "Returns 'OK' and item when the path matches",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			body:               "{ \"Item\": \"Apples\" }",
			output:             &data.Item{Item: "Apples", ID: "888"},
			outputErr:          nil,
			expectedRes:        &data.Item{Item: "Apples", ID: "888"},
			expectedStatusCode: 200,
			shouldCallDB:       true,
		},
		{
			name:               "Returns 'Bad Request' when the body is wrong/missing",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			body:               "",
			expectedRes:        nil,
			expectedStatusCode: 400,
			shouldCallDB:       false,
		},
		{
			name:               "Returns 'Bad Request' when the item does not exist",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			body:               "{ \"Item\": \"Apples\" }",
			output:             nil,
			outputErr:          db.NewError(db.ErrorNotFound),
			expectedRes:        nil,
			expectedStatusCode: 404,
			shouldCallDB:       true,
		},
		{
			name:               "Returns 'Internal Server Error' when the db returns an error",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			body:               "{ \"Item\": \"Apples\" }",
			output:             nil,
			outputErr:          errors.New("Something bad happened"),
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
				On("UpdateItem", &tt.listID, &tt.itemID, &tt.newName).
				Return(tt.output, tt.outputErr)

			d := patchItem{db: dbMocked}

			input := createAPIGatewayV2HTTPRequest(tt.path, "PATCH", tt.body)
			gotRes, statusCode := d.handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
