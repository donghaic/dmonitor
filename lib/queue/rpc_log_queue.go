package queue

import (
	"destroyer-monitor/config"
	"destroyer-monitor/lib/zap"
	"fmt"
	"github.com/johntech-o/gorpc"
	"github.com/serialx/hashring"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TcpLogQueue struct {
	lock      sync.RWMutex
	client    *gorpc.Client
	ring      *hashring.HashRing
	goodNodes []string
	badNodes  []string
}

func NewTcpLogQueue(cnf *config.Config) (*TcpLogQueue, error) {
	nodes := cnf.Logworker.Nodes
	if len(nodes) == 0 {
		panic("click log worker node must great than 1")
	}
	netOptions := gorpc.NewNetOptions(time.Second*10, time.Second*20, time.Second*20)
	weights := make(map[string]int)
	for _, node := range nodes {
		nodeAndWeight := strings.Split(node, "#")
		weight, err := strconv.Atoi(nodeAndWeight[1])
		if err != nil {
			weight = 1
		}
		weights[nodeAndWeight[0]] = weight
	}

	ring := hashring.New(nodes)
	client := gorpc.NewClient(netOptions)
	tcpLogQueue := &TcpLogQueue{client: client, ring: ring, goodNodes: nodes, badNodes: []string{}}
	go tcpLogQueue.pingTimer()
	return tcpLogQueue, nil
}

func (q *TcpLogQueue) Enqueue(clickId string, value []byte) error {
	var response string
	q.lock.RLock()
	node, ok := q.ring.GetNode(clickId)
	q.lock.RUnlock()
	if ok {
		err := q.client.CallWithAddress(node, "RpcService", "SaveClickLog", value, &response)
		if err != nil {
			zap.Get().Error("call rpc error data=", string(value))
		}
		return err
	} else {
		zap.Get().Error("enqueue error data=", string(value))
	}
	return nil
}

func (q *TcpLogQueue) Dequeue() (*Item, error) {
	panic("implement me")
}

func (q *TcpLogQueue) PeekByOffset(offset uint64) (*Item, error) {
	panic("implement me")
}
func (q *TcpLogQueue) pingTimer() {
	timer := time.NewTimer(time.Second * 2)
	for {
		select {
		case <-timer.C:
			q.ping()
			timer.Reset(time.Second * 2)
		}
	}
}

func (q *TcpLogQueue) ping() {
	q.lock.Lock()
	zap.Get().Info(fmt.Sprintf("start ping good:%d, bad:%d", len(q.goodNodes), len(q.badNodes)))
	var response string
	var goodNode []string
	for _, node := range q.goodNodes {
		err := q.client.CallWithAddress(node, "RpcService", "Ping", "ping", &response)
		if err != nil {
			q.ring = q.ring.RemoveNode(node)
			q.badNodes = append(q.badNodes, node)
			zap.Get().Error("can not connect to node=", node, err)
		} else {
			goodNode = append(goodNode, node)
		}
	}

	q.goodNodes = goodNode

	// 检测
	var badNode []string
	for _, node := range q.badNodes {
		err := q.client.CallWithAddress(node, "RpcService", "Ping", "ping", &response)
		if err == nil {
			zap.Get().Info("reconnect to node=", node)
			q.ring = q.ring.AddNode(node)
			q.goodNodes = append(q.goodNodes, node)
		} else {
			zap.Get().Error("can not connect to node=", node, err)
			badNode = append(badNode, node)
		}
	}
	q.badNodes = badNode

	q.lock.Unlock()
}
