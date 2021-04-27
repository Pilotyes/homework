//Package cache is a Wrappep for gp-cache with debug messages
package cache

import (
	"shop-api/internal/model"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

//Cache ...
type Cache struct {
	cache  *cache.Cache
	logger *logrus.Entry
}

//NewCache ...
func NewCache(logger *logrus.Entry, expirationTime, cleanupInterval time.Duration) *Cache {
	cache := cache.New(expirationTime*time.Second, cleanupInterval*time.Minute)

	return &Cache{
		cache:  cache,
		logger: logger,
	}
}

//Get ...
func (c *Cache) Get(id string) *model.Item {
	logger := c.logger.WithFields(logrus.Fields{
		"id": id,
	})

	logger.Debugln("Trying to get item from cache")
	if item, ok := c.cache.Get(id); ok {
		logger.Debugln("Item found in cache")
		return item.(*model.Item)
	}

	logger.Debugln("Item not found in cache")
	return nil
}

//Set ...
func (c *Cache) Set(id string, item *model.Item) {
	logger := c.logger.WithFields(logrus.Fields{
		"id": id,
	})

	c.cache.Set(id, item, cache.DefaultExpiration)
	logger.Debugln("Set item in cache")
}

//Delete ...
func (c *Cache) Delete(id string) {
	logger := c.logger.WithFields(logrus.Fields{
		"id": id,
	})

	c.Delete(id)
	logger.Debugln("Deleted item from cache")
}
