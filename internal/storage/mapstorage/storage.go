package mapstorage

import (
	"shop-api/internal/config"
	"shop-api/internal/model"
	"shop-api/internal/storage"
	"shop-api/internal/storage/internal/cache"

	"github.com/sirupsen/logrus"
)

//Storage ...
type Storage struct {
	itemsRepository *ItemsRepository
	config          *config.Config
}

//New ...
func New(config *config.Config) *Storage {
	return &Storage{
		config: config,
	}
}

//Items ...
func (s *Storage) Items() storage.ItemsRepository {
	if s.itemsRepository != nil {
		return s.itemsRepository
	}

	repositoryLogger := logrus.New()

	level, _ := logrus.ParseLevel(s.config.Server.LogLevel)
	repositoryLogger.SetLevel(level)

	logger := repositoryLogger.WithFields(logrus.Fields{
		"storage driver": config.InternalDriver,
	})

	cache := cache.NewCache(logger, s.config.Server.CacheExpirationTime, s.config.Server.CacheCleanupInterval)

	s.itemsRepository = &ItemsRepository{
		items:  make(map[model.ItemID]*model.Item),
		cache:  cache,
		logger: logger,
	}

	return s.itemsRepository
}
