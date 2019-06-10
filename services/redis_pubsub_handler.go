package services

import (
	"destroyer-monitor/dao"
	mredis "destroyer-monitor/lib/redis"
	"destroyer-monitor/lib/zap"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type RedisPubSubService struct {
	redisPool *mredis.ConnPool
	redisDao  *dao.RedisDao
}

func NewRedisPubSubService(pubsubPool *mredis.ConnPool, redisDao *dao.RedisDao) *RedisPubSubService {
	return &RedisPubSubService{pubsubPool, redisDao}
}

func (s *RedisPubSubService) Start() {
	zap.Get().Info("start subscrible redis")
	conn, e := s.redisPool.GetConnection()
	if e != nil {
		panic("can't sub to redis")
	}

	psConn := redis.PubSubConn{Conn: conn}
	_ = psConn.Subscribe("ps:customer:offer")
	_ = psConn.Subscribe("ps:blacklist")
	_ = psConn.Subscribe("ps:channel:info")

	for {
		data := ""
		switch msg := psConn.Receive().(type) {
		case redis.Message:
			switch msg.Channel {
			case "ps:customer:offer":
				data = string(msg.Data)
				ids := strings.Split(data, ":")
				if 0 < len(ids) {
					s.redisDao.Update("offer", ids[len(ids)-1])
				}

			case "ps:blacklist":
				data = string(msg.Data)
				ids := strings.Split(data, ":")
				if 0 < len(ids) {
					s.redisDao.Update("blacklist", ids[len(ids)-1])
				}

			case "ps:channel:info":
				data = string(msg.Data)
				ids := strings.Split(data, ":")
				if 0 < len(ids) {
					s.redisDao.Update("channel", ids[len(ids)-1])
				}
			}
		case error:
			zap.Get().Error("ReceiveMsg error: %v\n", msg)
			defer conn.Close()
			go s.Start()
			break
		}
	}
}
