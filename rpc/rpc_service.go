package rpc

import (
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/lib/zap"
)

const EchoContent = "pong"

type RpcService struct {
	queue queue.Queue
}

func NewRpcService(queue queue.Queue) RpcService {
	return RpcService{queue: queue}
}

func (r *RpcService) SaveClickLog(data []byte, res *string) error {
	err := r.queue.Enqueue("", data)
	if err != nil {
		zap.Get().Error("rpc enqueue error, data=", string(data))
	}
	res = nil
	return err
}

func (r *RpcService) Ping(req string, res *string) error {
	*res = EchoContent
	return nil
}
