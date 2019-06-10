package handler

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/models"
	"destroyer-monitor/services/delayed"
	"destroyer-monitor/services/handler/internal"
	"encoding/json"
	"fmt"
)

const (
	API_CONV      = "api"
	CUSTOMER_CONV = "customer"
)

type convHandler interface {
	Handle(params *models.EventParams, convLogEntity []models.EventLog, clickLogEntity *models.ClickLogEntity)
}

// 转化业务处理
type ConvHandler struct {
	convMongoDao     *dao.ConvMongoDao
	clickMongoDaoArr []dao.ClickMongoDao
	redisDao         *dao.RedisDao
	clickhouseDao    *dao.ClickhouseDao
	delayedTaskQ     *delayed.DelayedTaskQueue
	handlerMap       map[string]convHandler
}

func NewConvHandler(convMongoDao *dao.ConvMongoDao,
	clickMongoDaoArr []dao.ClickMongoDao,
	redisDao *dao.RedisDao,
	clickhouseDao *dao.ClickhouseDao,
	delayedTaskQ *delayed.DelayedTaskQueue) *ConvHandler {
	handlerMap := make(map[string]convHandler)
	handlerMap[API_CONV] = &internal.ApiConvHandler{ConvMongoDao: convMongoDao, ClickMongoDaoArr: clickMongoDaoArr, RedisDao: redisDao, ClickhouseDao: clickhouseDao}
	handlerMap[CUSTOMER_CONV] = &internal.CustomerConvHandler{ConvMongoDao: convMongoDao, ClickMongoDaoArr: clickMongoDaoArr, RedisDao: redisDao, ClickhouseDao: clickhouseDao}

	return &ConvHandler{
		convMongoDao:  convMongoDao,
		redisDao:      redisDao,
		clickhouseDao: clickhouseDao,
		delayedTaskQ:  delayedTaskQ,
		handlerMap:    handlerMap,
	}
}

// 核心业务处理逻辑
func (h *ConvHandler) Handle(params *models.EventParams) *models.Response {
	go logConvEvent(params)

	clkEvent, err := h.findClickEventById(params.ClickId, params.ClickTs)
	if err != nil {
		// 延迟处理转化
		_ = h.delayedTaskQ.Put(delayed.Task{Cnt: 0, Params: params})
		return &models.Response{Code: 200, Content: "OK"}
	}

	go h.doHandle(params, clkEvent)

	return &models.Response{Code: 200, Content: "OK"}

}

func (h *ConvHandler) doHandle(params *models.EventParams, clickEntity *models.ClickLogEntity) {
	convEventArr, err := h.convMongoDao.FindConvEventByClickIdAndOfferId(params.ClickId, params.OfferId)
	if nil != err {
		zap.Get().Error("event PostProcess get eventlog error.", err)
		return
	}
	offerType := clickEntity.OfferType
	if handler, ok := h.handlerMap[offerType]; ok {
		handler.Handle(params, convEventArr, clickEntity)
	} else {
		data, _ := json.Marshal(params)
		zap.Get().Error(fmt.Sprintf("unknown offer type %s, data=%s ", offerType, string(data)))
	}
}

func logConvEvent(event *models.EventParams) {
	data, _ := json.Marshal(event)
	line := fmt.Sprintf("conversion %s", string(data))
	zap.GetEvent().Info(line)
}

func (h *ConvHandler) findClickEventById(clickId string, clickTsStr string) (*models.ClickLogEntity, error) {
	var err error
	for _, clickDao := range h.clickMongoDaoArr {
		entity, e := clickDao.FindClickEventById(clickId, clickTsStr)
		err = e
		if entity != nil && 0 < len(entity.DestroyerId) {
			return entity, nil
		}
	}
	return nil, err
}
