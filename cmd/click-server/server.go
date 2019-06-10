package main

import (
	"destroyer-monitor/api"
	"destroyer-monitor/config"
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/lib/redis"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/services"
	"destroyer-monitor/services/handler"
	"destroyer-monitor/utils"
	"fmt"
	routing "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

type ClkServer struct {
	cnf          *config.Config
	clickApi     *api.ClickApi
	smartLinkApi *api.SmartLinkApi
	queue        *queue.Queue
}

func NewClickServer(cnf *config.Config) *ClkServer {
	return &ClkServer{cnf: cnf}
}

func (s *ClkServer) Init() error {
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

	logger.Info("start TcpLogQueue init")
	localQueue, err := queue.NewTcpLogQueue(s.cnf)
	if err != nil {
		logger.Error("NewLocalQueue init error", err)
		return err
	}

	logger.Info("start http client init")
	httpCli := utils.NewHttpCli(&s.cnf.Httpcli)

	clickHandler := handler.NewClickHandler(redisDao, localQueue, httpCli)
	s.clickApi = api.NewClick(clickHandler)
	s.smartLinkApi = &api.SmartLinkApi{}

	logger.Info("done server init")
	return nil

}

func (s *ClkServer) Run() {
	logger := zap.Get()

	router := routing.New()
	router.GET("/click", s.clickApi.Handle)
	router.POST("/smartlink", s.smartLinkApi.SetSmartLink)
	router.GET("/smartlink", s.smartLinkApi.GetSmartLink)
	router.GET("/clk/hello", hello)

	address := fmt.Sprintf(": %d", s.cnf.Port)
	logger.Info("ClkServer bind to ", address)
	if err := fasthttp.ListenAndServe(address, api.CORS(router.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
func hello(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Hello World!")
}
