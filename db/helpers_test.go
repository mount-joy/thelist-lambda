package db

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var re = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestGenerateID(t *testing.T) {
	t.Run("should return a UUIDv4", func(t *testing.T) {
		for i := 1; i <= 1000; i++ {
			id := generateID()

			t.Run(id+" is a valid UUIDv4", func(t *testing.T) {
				assert.Regexp(t, re, id)
			})
		}
	})

	t.Run("Shouldn't clash", func(t *testing.T) {
		setOfIDs := make(map[string]bool)
		n := 1000000

		for i := 1; i <= n; i++ {
			id := generateID()
			setOfIDs[id] = true
		}

		assert.Equal(t, n, len(setOfIDs))
	})
}
