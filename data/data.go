package data

// Item represents the data structure of an item on a list
type Item struct {
	ID     string `json:"Id"`
	ListID string `json:"ListId"`
	Item   string `json:"Item"`
}
