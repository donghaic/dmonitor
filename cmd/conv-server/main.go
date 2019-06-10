package main

import (
	"destroyer-monitor/api"
	"destroyer-monitor/config"
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/mongo"
	"destroyer-monitor/lib/redis"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/services"
	"destroyer-monitor/services/delayed"
	"destroyer-monitor/services/handler"
	"destroyer-monitor/utils"
	"fmt"
	routing "github.com/buaazp/fasthttprouter"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func main() {
	pflag.String("conf", "./config.yaml", "set configuration `file`")
	pflag.String("profile", "dev", "app profile")
	pflag.String("log-dir", "./logs", "server logs dir")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	configFile := viper.GetString("conf")
	var logger = zap.Get()
	if utils.IsEmpty(configFile) {
		pflag.Usage()
		logger.Info("can't not found config file")
		os.Exit(1)
	}

	cnf, err := config.ReadConfig(configFile)
	if err != nil {
		logger.Error("read cnf file error. file=", configFile)
		log.Fatalf("read cnf file error: %s", err)
	}

	server := &Server{cnf: cnf}
	_ = server.Init()
	server.Run()
}

type Server struct {
	cnf        *config.Config
	conversion *api.Conversion
}

func (s Server) Init() error {
	logger := zap.Get()
	logger.Info("start entity redis init")
	redisDao, err := dao.NewRedisDao(s.cnf)
	if err != nil {
		logger.Error("redisDao init error", err)
		return err
	}

	logger.Info("start pubsub redis init")
	pubsubPool, err := redis.NewPool(&s.cnf.Redis.Pubsub)
	if err != nil {
		logger.Error("pubsub redis init error", err)
		return err
	}

	logger.Info("start RedisPubSubService init")
	pubsubService := services.NewRedisPubSubService(pubsubPool, redisDao)
	pubsubService.Start()

	logger.Info("start task queue init dir=", s.cnf.Queue.LocalDataDir)
	delayedTaskQueue, err := delayed.NewDelayedQueue(s.cnf.Queue.LocalDataDir)
	if err != nil {
		logger.Error("delayed task queue init error", err)
		return err
	}

	convMgoPool, err := mongo.NewMgoPool(&s.cnf.Mongodb.Conversion)
	convMongoDao := dao.NewConvMongoDao(convMgoPool)

	click1MgoPool1, err := mongo.NewMgoPool(&s.cnf.Mongodb.Click1)
	click1MgoPoo2, err := mongo.NewMgoPool(&s.cnf.Mongodb.Click2)
	clickMongoDao1 := dao.NewClickMongoDao(click1MgoPool1)
	clickMongoDao2 := dao.NewClickMongoDao(click1MgoPoo2)
	clickMongoDaoArr := []dao.ClickMongoDao{*clickMongoDao1, *clickMongoDao2}

	// TODO clickhouse
	handler := handler.NewConvHandler(convMongoDao, clickMongoDaoArr, redisDao,
		nil, delayedTaskQueue)
	s.conversion = api.NewConversion(handler)
	return nil
}

func (s Server) Run() {
	logger := zap.Get()

	router := routing.New()
	router.GET("/e/hello/:name", hello)

	router.POST("/active", s.conversion.Handle)
	router.GET("/active", s.conversion.Handle)
	router.POST("/e", s.conversion.Handle)
	router.GET("/e", s.conversion.Handle)

	address := fmt.Sprintf(": %d", s.cnf.Port)
	logger.Info("ConvServer bind to ", address)
	if err := fasthttp.ListenAndServe(address, api.CORS(router.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
func hello(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Hello World!")
}
