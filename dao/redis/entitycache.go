package cache

import (
	"destroyer-monitor/dao/redis/internal"
	"destroyer-monitor/lib/redis"
)

type EntityCacheService struct {
	localCache *internal.LocalCacheService
	redisCache *internal.RedisCacheService
}

func NewEntityCacheService(redisPool *redis.ConnPool) *EntityCacheService {
	return &EntityCacheService{
		localCache: internal.NewLocalCache(),
		redisCache: internal.NewRedisCacheService(redisPool),
	}
}

