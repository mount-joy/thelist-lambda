package db

import (
	"github.com/mount-joy/thelist-lambda/data"
)

// DB - interface for talking to the database
type DB interface {
	DeleteItem(string, string) error
	GetItem(listID string, itemID string) (*data.Item, error)
	GetItemsOnList(string) (*[]data.Item, error)
	UpdateItem(string, string, string) (*data.Item, error)
}
