package test

import (
	"destroyer-monitor/models"
	"encoding/json"
	"flag"
	"github.com/johntech-o/gorpc"
	"testing"
	"time"
)

var client *gorpc.Client
var A string

func init() {
	netOptions := gorpc.NewNetOptions(time.Second*10, time.Second*20, time.Second*20)
	client = gorpc.NewClient(netOptions)
	flag.StringVar(&A, "a", "127.0.0.1:6668", "remote address:port")
	flag.Parse()
}

func TestRpcClient(t *testing.T) {

	var strReturn string

	_ = client.CallWithAddress(A, "RpcService", "Ping", "ping", &strReturn)
	println(strReturn)

	//call remote service SaveClickLog and method ClickLog
	clickLog := &models.ClickLog{}
	bytes, _ := json.Marshal(clickLog)
	err := client.CallWithAddress(A, "RpcService", "SaveClickLog", bytes, &strReturn)
	if err != nil {
		println(err.Error(), err.Errno())
	}
	println(strReturn)

}
