package main

import (
	"destroyer-monitor/config"
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/mongo"
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/rpc"
	"destroyer-monitor/services"
	"destroyer-monitor/utils"
	"fmt"
	"github.com/johntech-o/gorpc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	pflag.String("conf", "./config.yaml", "set configuration `file`")
	pflag.String("mongo", "click1", "点击保存到那个mongoDB click1 or click2")
	pflag.String("profile", "dev", "app profile")
	pflag.String("log-dir", "./logs", "server logs dir")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	configFile := viper.GetString("conf")
	mongoInstance := viper.GetString("mongo")
	var logger = zap.Get()
	if utils.IsEmpty(configFile) || utils.IsEmpty(mongoInstance) {
		pflag.Usage()
		logger.Info("can't not found config file")
		os.Exit(1)
	}

	cnf, err := config.ReadConfig(configFile)
	if err != nil {
		logger.Error("read cnf file error. file=", configFile)
		log.Fatalf("read cnf file error: %s", err)
	}

	startRpcServer(cnf, mongoInstance)

}

func startRpcServer(cnf *config.Config, mongoInstance string) {
	logger := zap.Get()
	logger.Info("start LocalQueue init")
	localQueue, err := queue.NewLocalQueue(cnf)
	if err != nil {
		logger.Error("NewLocalQueue init error", err)
		panic("localQueue init error")
	}

	logger.Info("start clickhouse dao init")
	clickhouseDao, err := dao.NewClickhouseDao(cnf)

	logger.Info("start mongo dao init")
	var option mongo.DBOption
	if utils.Equal(mongoInstance, "click1") {
		option = cnf.Mongodb.Click1
		zap.Get().Info("use click1 mongodb config")
	} else if utils.Equal(mongoInstance, "click2") {
		option = cnf.Mongodb.Click2
		zap.Get().Info("use click2 mongodb config")
	} else {
		panic(fmt.Sprintf("unknown mongo db %s", mongoInstance))
	}
	clickMgoPool, err := mongo.NewMgoPool(&option)
	clickMongoDao := dao.NewClickMongoDao(clickMgoPool)

	// 启动Worker线程处理点击转存任务
	clickWorker := services.NewClickWorker(localQueue, clickMongoDao, clickhouseDao)
	clickWorker.Run()

	// 启动RPC服务接收点击日志
	rpcService := rpc.NewRpcService(localQueue)
	s := gorpc.NewServer(cnf.Logworker.Address)
	_ = s.Register(&rpcService)
	logger.Info("start rpc server on address: ", cnf.Logworker.Address)
	s.Serve()
}
