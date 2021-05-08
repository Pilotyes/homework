package model

import (
	"testing"
)

func TestItem_IsEmpty(t *testing.T) {
	type fields struct {
		ID            ItemID
		Name          string
		Description   string
		OriginalPrice *float64
		DiscountPrice *float64
		Articul       int
		Category      string
		ProductOfDay  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Item is empty",
			fields: fields{},
			want:   true,
		},
		{
			name: "Item is not empty",
			fields: fields{
				Name:          "Item 1",
				Description:   "Description 1",
				OriginalPrice: new(float64),
				DiscountPrice: new(float64),
				Articul:       100001,
				Category:      "Category 1",
				ProductOfDay:  false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:            tt.fields.ID,
				Name:          tt.fields.Name,
				Description:   tt.fields.Description,
				OriginalPrice: tt.fields.OriginalPrice,
				DiscountPrice: tt.fields.DiscountPrice,
				Articul:       tt.fields.Articul,
				Category:      tt.fields.Category,
				ProductOfDay:  tt.fields.ProductOfDay,
			}
			if got := i.IsEmpty(); got != tt.want {
				t.Errorf("Item.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemID_GetString(t *testing.T) {
	tests := []struct {
		name string
		i    ItemID
		want string
	}{
		{
			name: "Valid ItemID",
			i:    ItemID(5),
			want: "5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetString(); got != tt.want {
				t.Errorf("ItemID.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}
