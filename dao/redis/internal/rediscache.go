package internal

import (
	"destroyer-monitor/lib/redis"
)

type RedisCacheService struct {
	pool *redis.ConnPool
}

func NewRedisCacheService(connPool *redis.ConnPool) *RedisCacheService {
	return &RedisCacheService{
		pool: connPool,
	}
}
