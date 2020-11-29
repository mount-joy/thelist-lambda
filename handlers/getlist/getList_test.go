package getlist

import (
	"errors"
	"testing"

	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/stretchr/testify/assert"
)

func TestGetListMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists/b6cf642d/",
			method:      "GET",
			expectedRes: true,
		},
		{
			name:        "Returns true for a uppercase ID",
			path:        "/lists/B6CF642D/",
			method:      "GET",
			expectedRes: true,
		},
		{
			name:        "Returns true without trailing slash",
			path:        "/lists/b6cf642d/",
			method:      "GET",
			expectedRes: true,
		},
		{
			name:        "Returns false for item path",
			path:        "/lists/b6cf642d/items/73bb82c4/",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false for list path without trailing slash",
			path:        "/lists/b6cf642ditems/73bb82c4",
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
			d := getList{db: dbMocked}
			gotRes := d.Match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

type mockGetList struct {
	res *data.List
	err error
}

func TestGetItemHandle(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		listID             string
		mockOutput         *mockGetList
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
			path:               "/lists",
			listID:             "test-list-id",
			expectedRes:        nil,
			expectedStatusCode: 500,
		},
		{
			name:               "Returns 'OK' and results when the path matches",
			path:               "/lists/test-list-id/",
			listID:             "test-list-id",
			mockOutput:         &mockGetList{res: &data.List{Name: "ABC", ListKey: data.ListKey{ID: "888"}}, err: nil},
			expectedRes:        &data.List{Name: "ABC", ListKey: data.ListKey{ID: "888"}},
			expectedStatusCode: 200,
		},
		{
			name:               "Returns 'Internal Server Error' when db returns an error",
			path:               "/lists/test-list-id/",
			listID:             "test-list-id",
			mockOutput:         &mockGetList{res: nil, err: errors.New("It went bad")},
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
					On("GetList", tt.listID).
					Return(tt.mockOutput.res, tt.mockOutput.err).
					Once()
			}

			d := getList{db: dbMocked}

			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, "GET", "")
			gotRes, statusCode := d.Handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
