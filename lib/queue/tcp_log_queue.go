package queue

import "destroyer-monitor/config"

type TcpLogQueue struct {
}

func NewTcpLogQueue(cnf *config.Config) (*TcpLogQueue, error) {
	return nil, nil
}

func (q *TcpLogQueue) Enqueue(clickId string, value []byte) error {
	return nil
}

func (q *TcpLogQueue) Dequeue() (*Item, error) {
	return nil, nil
}
