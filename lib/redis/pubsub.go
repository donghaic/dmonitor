package redis

import (
	predis "github.com/garyburd/redigo/redis"
)

type PubSubService struct {
	pool *ConnPool
}

func NewPubSub(redisPool *ConnPool) *PubSubService {
	return &PubSubService{pool: redisPool}
}

//
// Publish publish key value
func (s *PubSubService) Publish(key string, value string) error {
	conn, _ := s.pool.GetConnection()
	conn.Do("PUBLISH", key, value)
	return nil
}

// Subscribe subscribe
func (s *PubSubService) Subscribe(key string, callback func([]byte)) error {
	rc, _ := s.pool.GetConnection()
	psc := predis.PubSubConn{Conn: rc}
	if err := psc.Subscribe(key); err != nil {
		return err
	}

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case predis.Message:
				callback(v.Data)
			}
		}
	}()
	return nil
}
