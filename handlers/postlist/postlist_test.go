package postlist

import (
	"fmt"
	"testing"

	"github.com/mount-joy/thelist-lambda/data"
	"github.com/mount-joy/thelist-lambda/handlers/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestPostListMatch(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		method      string
		expectedRes bool
	}{
		{
			name:        "Returns true for a matching path",
			path:        "/lists",
			method:      "POST",
			expectedRes: true,
		},
		{
			name:        "Returns true with trailing slash",
			path:        "/lists/",
			method:      "POST",
			expectedRes: true,
		},
		{
			name:        "Returns false fot PUT request",
			path:        "/lists/",
			method:      "PUT",
			expectedRes: false,
		},
		{
			name:        "Returns false fot DELETE request",
			path:        "/lists/",
			method:      "DELETE",
			expectedRes: false,
		},
		{
			name:        "Returns false fot GET request",
			path:        "/lists/",
			method:      "GET",
			expectedRes: false,
		},
		{
			name:        "Returns false fot PATCH request",
			path:        "/lists/",
			method:      "PATCH",
			expectedRes: false,
		},
		{
			name:        "Returns false when path is empty",
			path:        "",
			method:      "POST",
			expectedRes: false,
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
			name:        "Returns false for item path",
			path:        "/lists/b6cf642d/items/",
			method:      "POST",
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := testhelpers.CreateAPIGatewayV2HTTPRequest(tt.path, tt.method, "")

			handler := &postList{}

			gotRes := handler.Match(input)

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

func TestPostListHandle(t *testing.T) {
	type mockPostList struct {
		res *data.List
		err error
	}
	tests := []struct {
		name               string
		listName           string
		mockPostList       *mockPostList
		badJsonInput       bool
		expectedRes        interface{}
		expectedStatusCode int
	}{
		{
			name:     "Returns 'OK' if database succeeds",
			listName: "myeList",
			mockPostList: &mockPostList{
				res: &data.List{Name: "myList", ListKey: data.ListKey{ID: "1234"}},
				err: nil,
			},
			expectedRes:        &data.List{Name: "myList", ListKey: data.ListKey{ID: "1234"}},
			expectedStatusCode: 200,
		},
		{
			name:     "Returns 'internal server error' if database fails",
			listName: "myeList",
			mockPostList: &mockPostList{
				res: nil,
				err: fmt.Errorf("uh oh"),
			},
			expectedRes:        nil,
			expectedStatusCode: 500,
		},
		{
			name:               "Returns error for bad json in body",
			badJsonInput:       true,
			expectedRes:        nil,
			expectedStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMocked := testhelpers.MockDB{}
			dbMocked.Test(t)
			defer dbMocked.AssertExpectations(t)

			if tt.mockPostList != nil {
				dbMocked.
					On("CreateList", tt.listName).
					Return(tt.mockPostList.res, tt.mockPostList.err).
					Once()
			}

			d := postList{db: &dbMocked}

			var body string
			if tt.badJsonInput {
				body = `badjson,`
			} else {
				body = fmt.Sprintf("{ \"Name\": %q }", tt.listName)
			}
			// Fine to hard code path and method as they aren't used in this function
			input := testhelpers.CreateAPIGatewayV2HTTPRequest("/lists/", "POST", body)

			gotRes, statusCode := d.Handle(input)

			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
