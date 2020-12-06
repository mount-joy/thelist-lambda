package patchitem

import (
	"errors"
	"testing"

	"github.com/mount-joy/thelist-lambda/db"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"
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
			dbMocked := &testhelpers.MockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, tt.method, "")
			d := patchItem{db: dbMocked}
			gotRes := d.Match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

type mockUpdateItem struct {
	res *data.Item
	err error
}

func TestPatchItemHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		itemID             string
		newName            string
		isCompleted        *bool
		body               string
		expectedRes        interface{}
		expectedStatusCode int
		mockOutput         *mockUpdateItem
	}{
		{
			name:               "Returns 'Bad Request' when the path is empty",
			path:               "",
			listID:             "test-list-id",
			expectedRes:        nil,
			expectedStatusCode: 400,
		},
		{
			name:               "Returns 'Bad Request' when the path is not in the correct format",
			path:               "/lists/test-list-id",
			listID:             "test-list-id",
			expectedRes:        nil,
			expectedStatusCode: 400,
		},
		{
			name:               "Returns 'OK' and item when the path matches",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			isCompleted:        testhelpers.BoolToPointer(false),
			body:               `{ "Name": "Apples", "IsCompleted": false }`,
			mockOutput:         &mockUpdateItem{res: &data.Item{Name: "Apples", ItemKey: data.ItemKey{ID: "888"}, IsCompleted: false}, err: nil},
			expectedRes:        &data.Item{Name: "Apples", IsCompleted: false, ItemKey: data.ItemKey{ID: "888"}},
			expectedStatusCode: 200,
		},
		{
			name:               "Returns 'OK' and item when IsCompleted is omitted",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			isCompleted:        nil,
			body:               `{ "Name": "Apples" }`,
			mockOutput:         &mockUpdateItem{res: &data.Item{Name: "Apples", ItemKey: data.ItemKey{ID: "888"}, IsCompleted: false}, err: nil},
			expectedRes:        &data.Item{Name: "Apples", IsCompleted: false, ItemKey: data.ItemKey{ID: "888"}},
			expectedStatusCode: 200,
		},
		{
			name:               "Returns 'OK' and item when the path matches when Name is omitted",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "",
			isCompleted:        testhelpers.BoolToPointer(true),
			body:               `{ "IsCompleted": true }`,
			mockOutput:         &mockUpdateItem{res: &data.Item{Name: "Bananas", ItemKey: data.ItemKey{ID: "888"}, IsCompleted: true}, err: nil},
			expectedRes:        &data.Item{Name: "Bananas", IsCompleted: true, ItemKey: data.ItemKey{ID: "888"}},
			expectedStatusCode: 200,
		},
		{
			name:               "Returns 'Bad Request' when name is empty",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "",
			isCompleted:        nil,
			body:               `{ "Name": "" }`,
			mockOutput:         &mockUpdateItem{res: nil, err: db.ErrorBadRequest},
			expectedRes:        nil,
			expectedStatusCode: 400,
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
		},
		{
			name:               "Returns 'Bad Request' when the item does not exist",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			body:               "{ \"Name\": \"Apples\" }",
			mockOutput:         &mockUpdateItem{res: nil, err: db.ErrorNotFound},
			expectedRes:        nil,
			expectedStatusCode: 404,
		},
		{
			name:               "Returns 'Internal Server Error' when the db returns an error",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			newName:            "Apples",
			body:               "{ \"Name\": \"Apples\" }",
			mockOutput:         &mockUpdateItem{res: nil, err: errors.New("Something bad happened")},
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
					On("UpdateItem", tt.listID, tt.itemID, tt.newName, tt.isCompleted).
					Return(tt.mockOutput.res, tt.mockOutput.err).
					Once()
			}

			d := patchItem{db: dbMocked}

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, "PATCH", tt.body)
			gotRes, statusCode := d.Handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
