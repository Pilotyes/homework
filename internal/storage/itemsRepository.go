package storage

import "shop-api/internal/model"

//ItemsRepository ...
type ItemsRepository interface {
	GetItems() []*model.Item
	GetItem(id model.ItemID) *model.Item
	PutItem(*model.Item) (*model.Item, error)
	DeleteItem(id model.ItemID) error
}
