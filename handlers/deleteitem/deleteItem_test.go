package deleteitem

import (
	"errors"
	"testing"

	"github.com/mount-joy/thelist-lambda/db"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestDeleteItemMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "DELETE",
			expectedRes: true,
		},
		{
			name:        "Returns true for a uppercase ID",
			path:        "/lists/B6CF642D/items/73BB82C4/",
			method:      "DELETE",
			expectedRes: true,
		},
		{
			name:        "Returns true without trailing slash",
			path:        "/lists/b6cf642d/items/73bb82c4",
			method:      "DELETE",
			expectedRes: true,
		},
		{
			name:        "Returns false for list path",
			path:        "/lists/b6cf642d/",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false for list path without trailing slash",
			path:        "/lists/b6cf642d",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path",
			path:        "/lists/",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false for lists path without trailing slash",
			path:        "/lists",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false for items path",
			path:        "/lists/b6cf642d/items/",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false when path is empty",
			path:        "",
			method:      "DELETE",
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
			d := deleteItem{db: dbMocked}
			gotRes := d.Match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

type mockDeleteItem struct {
	res *data.Item
	err error
}

func TestDeleteItemHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		itemID             string
		mockOutput         *mockDeleteItem
		expectedStatusCode int
	}{
		{
			name:               "Returns 'Bad Request' when the path is empty",
			path:               "",
			listID:             "test-list-id",
			expectedStatusCode: 400,
		},
		{
			name:               "Returns 'Bad Request' when the path is not in the correct format",
			path:               "/lists/test-list-id",
			listID:             "test-list-id",
			expectedStatusCode: 400,
		},
		{
			name:               "Returns 'OK' and item when the path matches",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			mockOutput:         &mockDeleteItem{res: &data.Item{Name: "Apples", ID: "888"}, err: nil},
			expectedStatusCode: 200,
		},
		{
			name:               "Returns 'Bad Request' when the item does not exist",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			mockOutput:         &mockDeleteItem{res: nil, err: db.ErrorNotFound},
			expectedStatusCode: 404,
		},
		{
			name:               "Returns 'Internal Server Error' when the db returns an error",
			path:               "/lists/test-list-id/items/test-item-id/",
			listID:             "test-list-id",
			itemID:             "test-item-id",
			mockOutput:         &mockDeleteItem{res: nil, err: errors.New("Something bad happened")},
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
					On("DeleteItem", tt.listID, tt.itemID).
					Return(tt.mockOutput.res, tt.mockOutput.err).
					Once()
			}

			d := deleteItem{db: dbMocked}

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, "DELETE", "")
			gotRes, statusCode := d.Handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, nil, gotRes)
		})
	}
}
