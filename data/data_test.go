package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNameFieldInJson(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedName string
		wantErr      bool
	}{
		{
			name:         "Given json with name in it as a string all works",
			body:         `{"Name": "bob"}`,
			expectedName: "bob",
			wantErr:      false,
		},
		{
			name:         "Given json with name in it as a int, error",
			body:         `{"Name": 74}`,
			expectedName: "",
			wantErr:      true,
		},
		{
			name:         "Given bad json, error",
			body:         `{"Name": "bob",}`,
			expectedName: "",
			wantErr:      true,
		},
		{
			name:         "Given extra fields, still works",
			body:         `{"Name": "bob", "age": 55}`,
			expectedName: "bob",
			wantErr:      false,
		},
		{
			name:         "no name field in json, return error",
			body:         `{"age": 55}`,
			expectedName: "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotErr := GetNameFieldInJson(tt.body)

			assert.Equal(t, tt.expectedName, gotName)
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
