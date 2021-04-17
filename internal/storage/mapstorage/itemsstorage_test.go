package mapstorage

import (
	"reflect"
	"shop-api/internal/model"
	"testing"
)

func TestItemsRepository_GetItems(t *testing.T) {
	type fields struct {
		items map[int]*model.Item
	}
	tests := []struct {
		name   string
		fields fields
		want   []*model.Item
	}{
		{
			name: "Valid",
			fields: fields{
				items: map[int]*model.Item{
					1: {
						Id:          1,
						Name:        "Item 1",
						Description: "Desription 1",
						Price:       1.0,
					},
					2: {
						Id:          2,
						Name:        "Item 2",
						Description: "Desription 2",
						Price:       2.0,
					},
				},
			},
			want: []*model.Item{
				{
					Id:          1,
					Name:        "Item 1",
					Description: "Desription 1",
					Price:       1.0,
				},
				{
					Id:          2,
					Name:        "Item 2",
					Description: "Desription 2",
					Price:       2.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ItemsRepository{
				items: tt.fields.items,
			}
			if got := i.GetItems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsRepository.GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemsRepository_PutItem(t *testing.T) {
	type fields struct {
		items map[int]*model.Item
	}
	type args struct {
		item *model.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Item
		wantErr bool
	}{
		{
			name: "Valid",
			fields: fields{
				items: make(map[int]*model.Item),
			},
			args: args{
				item: &model.Item{
					Name:        "Item 1",
					Description: "Desription 1",
					Price:       1.0,
				},
			},
			want: &model.Item{
				Id:          1,
				Name:        "Item 1",
				Description: "Desription 1",
				Price:       1.0,
			},
			wantErr: false,
		},
		{
			name: "Nil item",
			fields: fields{
				items: make(map[int]*model.Item),
			},
			args: args{
				item: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ItemsRepository{
				items: tt.fields.items,
			}
			got, err := i.PutItem(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemsRepository.PutItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsRepository.PutItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemsRepository_DeleteItem(t *testing.T) {
	type fields struct {
		items map[int]*model.Item
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid",
			fields: fields{
				items: map[int]*model.Item{
					1: {
						Name:        "Item 1",
						Description: "Description 1",
						Price:       1.0,
					},
				},
			},
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "ID not found",
			fields: fields{
				items: map[int]*model.Item{},
			},
			args: args{
				id: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ItemsRepository{
				items: tt.fields.items,
			}
			if err := i.DeleteItem(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ItemsRepository.DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
