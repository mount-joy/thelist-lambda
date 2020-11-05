package data

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
