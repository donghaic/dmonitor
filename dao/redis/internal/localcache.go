package internal

import (
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var (
	defaultExpiration = 5 * time.Minute
	cleanupInterval   = 5 * time.Minute
)

type LocalCacheService struct {
	creative *gocache.Cache
	link     *gocache.Cache
}

func NewLocalCache() *LocalCacheService {
	return &LocalCacheService{
		creative: gocache.New(defaultExpiration, cleanupInterval),
		link:     gocache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *LocalCacheService) Flush() {
	c.link.Flush()
	c.creative.Flush()
}
