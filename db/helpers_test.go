package db

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	re := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	t.Run("should return a UUIDv4", func(t *testing.T) {
		for i := 1; i <= 1000; i++ {
			id := generateID()

			t.Run(id+" is a valid UUIDv4", func(t *testing.T) {
				assert.Regexp(t, re, id)
			})
		}
	})

	t.Run("Shouldn't clash", func(t *testing.T) {
		n := 1000000
		setOfIDs := make(map[string]bool, n)

		for i := 1; i <= n; i++ {
			id := generateID()
			setOfIDs[id] = true
		}

		assert.Equal(t, n, len(setOfIDs))
	})
}

func TestGetTimestamp(t *testing.T) {
	n := 100
	setOfTimestamps := make(map[string]bool, n)

	for i := 1; i <= n; i++ {
		id := getTimestamp()
		setOfTimestamps[id] = true
		time.Sleep(1 * time.Microsecond) // ensure some time has passed between each one
	}

	t.Run("Timestamps shouldn't clash", func(t *testing.T) {
		assert.Equal(t, n, len(setOfTimestamps))
	})

	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.\d+)?Z$`)

	t.Run("should return a valid ISO8601 date", func(t *testing.T) {
		for timestamp := range setOfTimestamps {
			t.Run(timestamp+" is a valid ISO8601 date", func(t *testing.T) {
				assert.Regexp(t, re, timestamp)
			})
		}
	})
}
