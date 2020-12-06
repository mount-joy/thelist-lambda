package data

import (
	"encoding/json"
	"fmt"
)

// ListKey represents the primary key of a list
type ListKey struct {
	ID string `json:"Id"`
}

// List represents the data structure of a list
type List struct {
	ListKey
	Name             string `json:"Name"`
	CreatedTimestamp string `json:"Created"`
	UpdatedTimestamp string `json:"Updated"`
}

// ItemKey represents the primary key of an item
type ItemKey struct {
	ID     string `json:"Id"`
	ListID string `json:"ListId"`
}

// Item represents the data structure of an item on a list
type Item struct {
	ItemKey
	Name             string `json:"Name"`
	IsCompleted      bool   `json:"IsCompleted"`
	CreatedTimestamp string `json:"Created"`
	UpdatedTimestamp string `json:"Updated"`
}

// GetNameFieldInJson gets the value of "Name" from the passed in json
func GetNameFieldInJson(jsonInput string) (string, error) {
	type PostInput struct {
		Name string `json:"Name"`
	}

	var input PostInput
	err := json.Unmarshal([]byte(jsonInput), &input)
	if err != nil {
		return "", err
	}

	if input.Name == "" {
		return "", fmt.Errorf("No \"Name\" field in the json")
	}

	return input.Name, err
}
