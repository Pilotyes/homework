package mapstorage

import (
	"shop-api/internal/model"
	"shop-api/internal/storage"
)

type Storage struct {
	itemsRepository *ItemsRepository
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Items() storage.ItemsRepository {
	if s.itemsRepository != nil {
		return s.itemsRepository
	}

	s.itemsRepository = &ItemsRepository{
		items: make(map[int]*model.Item),
	}

	return s.itemsRepository
}
