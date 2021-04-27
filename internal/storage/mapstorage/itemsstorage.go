package mapstorage

import (
	"errors"
	"shop-api/internal/model"
	"sync"

	"shop-api/internal/storage/internal/cache"

	"github.com/sirupsen/logrus"
)

//ItemsRepository ...
type ItemsRepository struct {
	nextID model.ItemID
	mutex  sync.Mutex
	items  map[model.ItemID]*model.Item
	cache  *cache.Cache
	logger *logrus.Entry
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
func (i *ItemsRepository) GetItem(id model.ItemID) *model.Item {
	/*
		if item := i.cache.Get(id.GetString()); item != nil {
			return item
		}
	*/

	return i.items[id]
}

//PutItem ...
func (i *ItemsRepository) PutItem(item *model.Item) (*model.Item, error) {
	if item == nil {
		return nil, errors.New("Item is null")
	}

	item.ID = model.ItemID(i.nextID)

	logger := i.logger.WithFields(logrus.Fields{
		"id": item.ID,
	})

	i.mutex.Lock()
	i.items[item.ID] = item
	i.mutex.Unlock()
	logger.Infof("Created new item: %#v\n", item)

	/*
		i.cache.Set(item.ID.GetString(), item)
	*/

	i.nextID++
	return item, nil
}

//DeleteItem ...
func (i *ItemsRepository) DeleteItem(id model.ItemID) error {
	if _, ok := i.items[id]; !ok {
		return errors.New("ID not found")
	}

	delete(i.items, id)
	/*
		i.cache.Delete(id.GetString())
	*/

	return nil
}
