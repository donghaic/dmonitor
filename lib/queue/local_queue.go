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
	err := os.MkdirAll(cnf.Queue.LocalDataDir, os.ModePerm)
	if err != nil {
		zap.Get().Info("", err)
	}
	q, err := goque.OpenQueue(cnf.Queue.LocalDataDir)
	if err != nil {
		return nil, err
	}
	return &LocalQueue{q,}, nil
}

func (q *LocalQueue) Enqueue(clickId string, value []byte) error {
	_, e := q.q.Enqueue(value)
	return e
}

func (q *LocalQueue) Dequeue() (*Item, error) {
	item, e := q.q.Dequeue()
	if e != nil {
		return nil, e
	}
	return &Item{item.Key, item.Value}, nil
}

func (q *LocalQueue) PeekByOffset(offset uint64) (*Item, error) {
	item, e := q.q.PeekByOffset(offset)
	if e != nil {
		return nil, e
	}
	return &Item{item.Key, item.Value}, nil
}
