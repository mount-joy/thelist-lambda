package data

import (
	"encoding/json"
	"fmt"
)

// List represents the data structure of a list
type List struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// Item represents the data structure of an item on a list
type Item struct {
	ID     string `json:"Id"`
	ListID string `json:"ListId"`
	Name   string `json:"Name"`
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
