package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type ConnPool struct {
	raw *redis.Pool
}

func NewPool(opt *PoolOption) (*ConnPool, error) {
	rawPool := newPool(opt)
	if rawPool == nil {
		return nil, errors.New("init redis raw failed")
	}
	// 测试连接池是否正常
	conn := rawPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil || conn.Err() != nil {
		return nil, conn.Err()
	}
	return &ConnPool{rawPool}, nil
}

func newPool(opt *PoolOption) *redis.Pool {
	return &redis.Pool{
		MaxActive:   opt.MaxActive,
		MaxIdle:     opt.MaxIdle,
		IdleTimeout: time.Duration(opt.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", opt.Address, redis.DialDatabase(opt.DB))
			if err != nil {
				return nil, err
			}
			if opt.Password != "" {
				if _, err := c.Do("AUTH", opt.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

/*** Single KV invoke ***/
func (p *ConnPool) SetString(key string, value interface{}) (interface{}, error) {
	conn := p.raw.Get()
	defer conn.Close()
	return conn.Do("SET", key, value)
}

func (p *ConnPool) GetString(key string) (string, error) {
	conn := p.raw.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (p *ConnPool) DelKey(key string) (interface{}, error) {
	conn := p.raw.Get()
	defer conn.Close()
	return conn.Do("DEL", key)
}

func (p *ConnPool) ExpireKey(key string, expireTime int64) (interface{}, error) {
	conn := p.raw.Get()
	defer conn.Close()
	return conn.Do("EXPIRE", key, expireTime)
}

func (p *ConnPool) GetConnection() (redis.Conn, error) {
	conn := p.raw.Get()
	return conn, conn.Err()
}
