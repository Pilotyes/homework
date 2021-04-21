package model

import "strconv"

//ItemID ...
type ItemID int

//Item ...
type Item struct {
	ID          ItemID   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
}

//IsEmpty ...
func (i *Item) IsEmpty() bool {
	return i.Name == "" && i.Description == "" && i.Price == nil
}

//GetString converts int to string
func (i ItemID) GetString() string {
	return strconv.Itoa(int(i))
}
