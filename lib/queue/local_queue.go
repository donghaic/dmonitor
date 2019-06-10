package queue

import (
	"destroyer-monitor/config"
	"destroyer-monitor/lib/zap"
	"github.com/beeker1121/goque"
	"os"
)

type LocalQueue struct {
	q *goque.Queue
}

func NewLocalQueue(cnf *config.Config) (*LocalQueue, error) {
	err := os.MkdirAll(cnf.LocalQueueDataDir, os.ModePerm)
	if err != nil {
		zap.Get().Info("", err)
	}
	q, err := goque.OpenQueue(cnf.LocalQueueDataDir)
	if err != nil {
		return nil, err
	}
	return &LocalQueue{q,}, nil
}

func (l *LocalQueue) Enqueue(clickId string, value []byte) error {
	_, e := l.q.Enqueue(value)
	return e
}

func (l *LocalQueue) Dequeue() (*Item, error) {
	item, e := l.q.Dequeue()
	if e != nil {
		return nil, e
	}
	return &Item{item.Key, item.Value}, nil
}
