package mapstorage

import (
	"reflect"
	"shop-api/internal/model"
	"shop-api/internal/storage/internal/cache"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestItemsRepository_GetItems(t *testing.T) {
	price := 1.0

	type fields struct {
		items map[model.ItemID]*model.Item
	}
	tests := []struct {
		name   string
		fields fields
		want   []*model.Item
	}{
		{
			name: "Valid",
			fields: fields{
				items: map[model.ItemID]*model.Item{
					1: {
						ID:            1,
						Name:          "Item 1",
						Description:   "Desription 1",
						OriginalPrice: &price,
						DiscountPrice: &price,
						Articul:       100001,
						Category:      "Category 1",
						ProductOfDay:  false,
					},
					2: {
						ID:            2,
						Name:          "Item 2",
						Description:   "Desription 2",
						OriginalPrice: &price,
						DiscountPrice: &price,
						Articul:       100002,
						Category:      "Category 2",
						ProductOfDay:  false,
					},
				},
			},
			want: []*model.Item{
				{
					ID:            1,
					Name:          "Item 1",
					Description:   "Desription 1",
					OriginalPrice: &price,
					DiscountPrice: &price,
					Articul:       100001,
					Category:      "Category 1",
					ProductOfDay:  false,
				},
				{
					ID:            2,
					Name:          "Item 2",
					Description:   "Desription 2",
					OriginalPrice: &price,
					DiscountPrice: &price,
					Articul:       100002,
					Category:      "Category 2",
					ProductOfDay:  false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logrus.New().WithFields(nil)
			i := &ItemsRepository{
				items: tt.fields.items,
				cache: cache.NewCache(logger, time.Duration(0), time.Duration(0)),
			}
			if got := i.GetItems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemsRepository.GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemsRepository_PutItem(t *testing.T) {
	price := 1.0

	type fields struct {
		items map[model.ItemID]*model.Item
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
				items: make(map[model.ItemID]*model.Item),
			},
			args: args{
				item: &model.Item{
					Name:          "Item 1",
					Description:   "Desription 1",
					OriginalPrice: &price,
					DiscountPrice: &price,
					Articul:       100001,
					Category:      "Category 1",
					ProductOfDay:  false,
				},
			},
			want: &model.Item{
				ID:            1,
				Name:          "Item 1",
				Description:   "Desription 1",
				OriginalPrice: &price,
				DiscountPrice: &price,
				Articul:       100001,
				Category:      "Category 1",
				ProductOfDay:  false,
			},
			wantErr: false,
		},
		{
			name: "Nil item",
			fields: fields{
				items: make(map[model.ItemID]*model.Item),
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
			logger := logrus.New().WithFields(nil)
			i := &ItemsRepository{
				nextID: 1,
				items:  tt.fields.items,
				cache:  cache.NewCache(logger, time.Duration(0), time.Duration(0)),
				logger: logrus.New().WithFields(nil),
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
	price := 1.0

	type fields struct {
		items map[model.ItemID]*model.Item
	}
	type args struct {
		id model.ItemID
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
				items: map[model.ItemID]*model.Item{
					1: {
						Name:          "Item 1",
						Description:   "Description 1",
						OriginalPrice: &price,
						DiscountPrice: &price,
						Articul:       100001,
						Category:      "Category 1",
						ProductOfDay:  false,
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
				items: map[model.ItemID]*model.Item{},
			},
			args: args{
				id: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logrus.New().WithFields(nil)
			i := &ItemsRepository{
				items: tt.fields.items,
				cache: cache.NewCache(logger, time.Duration(0), time.Duration(0)),
			}
			if err := i.DeleteItem(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ItemsRepository.DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
