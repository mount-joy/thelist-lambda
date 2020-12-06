package db

import (
	"time"

	"github.com/google/uuid"
)

func generateID() string {
	return uuid.New().String()
}

func getTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
