package model

import (
	"testing"
)

func TestItem_IsEmpty(t *testing.T) {
	type fields struct {
		ID          ItemID
		Name        string
		Description string
		Price       *float64
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
				Name:  "Item 1",
				Price: new(float64),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Item{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Price:       tt.fields.Price,
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
