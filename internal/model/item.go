package model

import "strconv"

//ItemID ...
type ItemID int

//Item ...
type Item struct {
	ID            ItemID   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	OriginalPrice *float64 `json:"original_price,omitempty"`
	DiscountPrice *float64 `json:"discount_price,omitempty"`
	Articul       int      `json:"articul,omitempty"`
	Category      string   `json:"category,omitempty"`
	ProductOfDay  bool     `json:"product_of_day,omitempty"`
}

//IsEmpty ...
func (i *Item) IsEmpty() bool {
	return i.Name == "" || i.Description == "" || i.OriginalPrice == nil ||
		i.DiscountPrice == nil || i.Articul == 0 || i.Category == ""
}

//GetString converts int to string
func (i ItemID) GetString() string {
	return strconv.Itoa(int(i))
}
