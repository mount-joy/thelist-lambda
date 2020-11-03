package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConf(t *testing.T) {
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
			confMocked := &conf{
				getEnv: func(key string) string {
					if key == "ENV" {
						return tt.runtimeEnv
					}
					return fmt.Sprintf("env_%s", key)
				},
			}

			gotRes := confMocked.getConf()

			assert.Equal(t, tt.expectedRes, gotRes)
		})
	}
}

func TestGetConfiguration(t *testing.T) {
	conf := GetConfiguration()

	assert.Greater(t, len(conf.Endpoint), 0)
	assert.Greater(t, len(conf.TableNames.Items), 0)
	assert.Greater(t, len(conf.TableNames.Lists), 0)
}
