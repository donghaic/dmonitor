package test

import (
	"destroyer-monitor/models"
	"encoding/json"
	"flag"
	"github.com/johntech-o/gorpc"
	"testing"
)

type RpcService struct {
}

func (r *RpcService) Send(clickLog *models.ClickLog, res *string) error {
	data, _ := json.Marshal(clickLog)
	println(string(data))
	return nil
}

var L string

func init() {
	flag.StringVar(&L, "l", "127.0.0.1:6668", "remote address:port")
	flag.Parse()
}

func TestRpcServer(t *testing.T) {
	s := gorpc.NewServer(L)
	s.Register(new(RpcService))
	s.Serve()
	panic("server fail")
}
