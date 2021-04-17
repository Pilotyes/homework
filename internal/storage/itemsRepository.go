package storage

import "shop-api/internal/model"

type ItemsRepository interface {
	GetItems() []*model.Item
	GetItem(id int) *model.Item
	PutItem(*model.Item) (*model.Item, error)
	DeleteItem(id int) error
}
