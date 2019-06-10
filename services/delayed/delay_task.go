package delayed

import (
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/models"
	"destroyer-monitor/utils"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"strconv"
	"time"
)

var Task_Queue_Bucket_Name = []byte("task_queue")
var Task_Cnt_Bucket_Name = []byte("task_queue_count")

type DelayedTaskQueue struct {
	db           *bolt.DB
	EventChannel chan Task
}

type Task struct {
	Cnt    int
	Params *models.EventParams
}

type Worker interface {
	Run(task Task)
}

func NewDelayedQueue(dataFile string) (*DelayedTaskQueue, error) {
	_ = os.MkdirAll(dataFile, os.ModePerm)

	db, err := bolt.Open(fmt.Sprintf("%s/task.db", dataFile), 0600, nil)
	if err != nil {
		return nil, err
	}
	err1 := db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists(Task_Queue_Bucket_Name)
		_, _ = tx.CreateBucketIfNotExists(Task_Cnt_Bucket_Name)
		return err
	})
	if err1 != nil {
		return nil, err1
	}
	ch := make(chan Task, 20)
	taskQueue := &DelayedTaskQueue{db, ch,}
	go taskQueue.startTimer()
	return taskQueue, nil
}

func (d *DelayedTaskQueue) Put(task Task) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(Task_Queue_Bucket_Name)
		rawKey := []byte(task.Params.ClickId)
		rawValue, _ := json.Marshal(task.Params)
		err := b.Put(rawKey, rawValue)

		bc := tx.Bucket(Task_Cnt_Bucket_Name)
		cntKey := fmt.Sprintf("cnt:%s", task.Params.ClickId)
		cntVal := string(task.Cnt)
		_ = bc.Put([]byte(cntKey), []byte(cntVal))
		return err
	})
}

func (d *DelayedTaskQueue) Delete(impId string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(Task_Queue_Bucket_Name)
		rawKey := []byte(impId)
		err := b.Delete(rawKey)

		bc := tx.Bucket(Task_Cnt_Bucket_Name)
		cntKey := fmt.Sprintf("cnt:%s", impId)
		_ = bc.Delete([]byte(cntKey))
		return err
	})
}

func (d *DelayedTaskQueue) startTimer() {
	zap.Get().Info("start delay task queue timer")
	loop := time.NewTimer(time.Minute * 1)
	for {
		select {
		case <-loop.C:
			_ = d.db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket(Task_Queue_Bucket_Name)
				bc := tx.Bucket(Task_Cnt_Bucket_Name)
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					cntKey := fmt.Sprintf("cnt:%s", utils.ByteToStr(k))
					reqParams := models.EventParams{}
					cnt, _ := strconv.Atoi(string(bc.Get([]byte(cntKey))))
					_ = json.Unmarshal(v, &reqParams)
					d.EventChannel <- Task{
						Cnt:    cnt,
						Params: &reqParams,
					}
				}
				return nil
			})
			loop.Reset(time.Minute * 1)
		}
	}
}
