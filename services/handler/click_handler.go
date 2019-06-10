package handler

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/models"
	"destroyer-monitor/services/handler/internal"
	"destroyer-monitor/utils"
	"encoding/json"
)

const (
	API_CLICK      = "1"
	CUSTOMER_CLICK = "2"
)

type ClickHandler struct {
	handlerMap map[string]handler
}

func NewClickHandler(redisDao *dao.RedisDao, queue queue.Queue, httpCli *utils.HttpConPool) *ClickHandler {
	handlerMap := make(map[string]handler)
	handlerMap[API_CLICK] = &internal.ApiClickHandler{RedisDao: redisDao, Queue: queue}
	handlerMap[CUSTOMER_CLICK] = &internal.CustomerClickHandler{RedisDao: redisDao, Queue: queue, HttpCli: httpCli}
	return &ClickHandler{handlerMap}
}

type handler interface {
	Handle(*models.ClickParams) *models.Response
}

func (c *ClickHandler) Handle(params *models.ClickParams) *models.Response {
	api := params.Api
	handler, ok := c.handlerMap[api]
	if ok {
		return handler.Handle(params)
	} else {
		data, _ := json.Marshal(params)
		zap.Get().Error("bad request, data=", string(data))
		return &models.Response{
			Code:    200,
			Content: "parameter error [api]",
		}
	}
}
