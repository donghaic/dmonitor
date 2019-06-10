package services

import (
	"destroyer-monitor/lib/queue"
	"time"
)

type ClickToMongoWorker struct {
	lqueue queue.Queue
}

func NewClickWorker(lqueue queue.Queue) *ClickToMongoWorker {
	return &ClickToMongoWorker{lqueue}
}

func (w *ClickToMongoWorker) Run(worker func(*queue.Item)) {
	if worker == nil {
		panic("worker is nil, please register a worker function first")
	}

	for {
		item, e := w.lqueue.Dequeue()
		if item == nil || e != nil {
			time.Sleep(100)
		}
		go worker(item)
	}
}
