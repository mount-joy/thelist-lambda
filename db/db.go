package db

import (
	"github.com/mount-joy/thelist-lambda/data"
)

// DB - interface for talking to the database
type DB interface {
	GetItemsOnList(*string) (*[]data.Item, error)
}
