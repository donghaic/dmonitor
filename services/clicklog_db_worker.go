package services

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/models"
	"encoding/json"
	"fmt"
	"time"
)

type ClickToDBWorker struct {
	q             queue.Queue
	mongoDao      *dao.ClickMongoDao
	clickhouseDao *dao.ClickhouseDao
}

func NewClickWorker(q queue.Queue, mongoDao *dao.ClickMongoDao, clickhouseDao *dao.ClickhouseDao) *ClickToDBWorker {
	return &ClickToDBWorker{q, mongoDao, clickhouseDao}
}

func (w *ClickToDBWorker) Run() {
	go func() {
		zap.Get().Info("start run click log to db worker loop")

		for {
			item, e := w.q.Dequeue()
			if item == nil || e != nil {
				time.Sleep(300)
			} else {
				go w.doWork(item)
			}
		}
	}()
}

func (w *ClickToDBWorker) doWork(item *queue.Item) {
	clickLog := models.ClickLog{}
	err := json.Unmarshal(item.Value, clickLog)
	if err != nil {
		zap.Get().Error(fmt.Sprintf("parse click log json error. data=%s", string(item.Value)), err)
	}

	println(string(item.Value))

	// TODO 转存
	// mongo
	// clickhouse
}
