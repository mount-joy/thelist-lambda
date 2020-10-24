package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		runtimeEnv  string
		expectedRes Config
	}{
		{
			name:       "When environment is dev then hardcoded values are used",
			runtimeEnv: "DEV",
			expectedRes: Config{
				Endpoint: "http://localhost:8000",
				TableNames: TableNames{
					Items: "items",
					Lists: "lists",
				},
			},
		},
		{
			name:       "When environment is prod then environment variables are used",
			runtimeEnv: "PROD",
			expectedRes: Config{
				Endpoint: "",
				TableNames: TableNames{
					Items: "env_TABLE_NAME_ITEMS",
					Lists: "env_TABLE_NAME_LISTS",
				},
			},
		},
		{
			name:       "When environment is nonsense then fallsback to dev values",
			runtimeEnv: "nonsense",
			expectedRes: Config{
				Endpoint: "http://localhost:8000",
				TableNames: TableNames{
					Items: "items",
					Lists: "lists",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getEnv := func(key string) string {
				if key == "ENV" {
					return tt.runtimeEnv
				}
				return fmt.Sprintf("env_%s", key)
			}

			confMocked := &conf{
				getEnv: getEnv,
			}

			gotRes := confMocked.getConf()

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}
