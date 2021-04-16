package mapstorage

import (
	"errors"
	"shop-api/internal/model"
)

type ItemsRepository struct {
	items map[int]*model.Item
}

//GetItems ...
func (i *ItemsRepository) GetItems() []*model.Item {
	var result []*model.Item = make([]*model.Item, 0, len(i.items))

	for id := range i.items {
		result = append(result, i.GetItem(id))
	}

	return result
}

//GetItem ...
func (i *ItemsRepository) GetItem(id int) *model.Item {
	return i.items[id]
}

//PutItem ...
func (i *ItemsRepository) PutItem(item *model.Item) (*model.Item, error) {
	if item == nil {
		return nil, errors.New("Item is null")
	}

	newID := len(i.items) + 1
	item.Id = newID
	i.items[newID] = item

	return item, nil
}

//DeleteItem ...
func (i *ItemsRepository) DeleteItem(id int) error {
	if _, ok := i.items[id]; !ok {
		return errors.New("ID not found")
	}

	delete(i.items, id)

	return nil
}
