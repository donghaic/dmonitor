package handler

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/models"
	"destroyer-monitor/services/delayed"
	"destroyer-monitor/services/macro"
	"destroyer-monitor/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 转化业务处理
type ConvHandler struct {
	mongoDao      *dao.MongoDao
	redisDao      *dao.RedisDao
	clickhouseDao *dao.ClickhouseDao
	delayedTaskQ  *delayed.DelayedTaskQueue
}

// 核心业务处理逻辑
func (h *ConvHandler) Handle(params *models.EventParams) *models.Response {
	go logConvEvent(params)

	clkEvent, err := h.mongoDao.FindClickEventById(params.ClickId, params.ClickTs)
	if err != nil {
		// 延迟处理转化
		_ = h.delayedTaskQ.Put(delayed.Task{Cnt: 0, Params: params})
		return &models.Response{Code: 200, Content: "OK"}
	}

	go h.doHandle(params, clkEvent)

	return &models.Response{Code: 200, Content: "OK"}

}

func (h *ConvHandler) doHandle(params *models.EventParams, clickEntity *models.ClickLogEntity) {
	convEventArr, err := h.mongoDao.FindConvEventByClickIdAndOfferId(params.ClickId, params.OfferId)
	if nil != err {
		zap.Get().Error("event PostProcess get eventlog error.", err)
		return
	}

	if utils.Equal(clickEntity.OfferType, "api") {
		h.handleApi(params, convEventArr, clickEntity)
	} else if utils.Equal(clickEntity.OfferType, "customer") {
		h.handleCustomer(params, convEventArr, clickEntity)
	} else {
		zap.Get().Error("unknown offer type", clickEntity.OfferType)
	}
}

func logConvEvent(event *models.EventParams) {
	data, _ := json.Marshal(event)
	line := fmt.Sprintf("conversion %s", string(data))
	zap.GetEvent().Info(line)
}

func (h *ConvHandler) handleApi(params *models.EventParams, convLogEntity []models.EventLog, clickLogEntity *models.ClickLogEntity) {
	channel, _ := h.redisDao.GetChannel(strconv.Itoa(clickLogEntity.ChannelId))
	//judget event type
	params.EventType = ""
	if 0 < len(convLogEntity) {
		log.Println("handleApi. event logs length gt 1: repeat")
		params.EventType = "repeat"
	}

	//can not find the channel info
	if 0 == channel.Id && 0 == len(params.EventType) {
		log.Println("handleApi. channel info not found: test")
		params.EventType = "test"
	}

	offcvs, errLog := h.redisDao.GetCampaignDayStats(false, params.OfferId)
	if nil != errLog {
		log.Println("handleApi. offer stats day error.", errLog)
		//错误日志处理
	}

	if offcvs > 0 && offcvs >= clickLogEntity.OfferCap && 0 == len(params.EventType) {
		log.Println("handleApi. offer stats day reach: test", offcvs, clickLogEntity.OfferCap)
		params.EventType = "test"
	}

	if 0 == len(params.EventType) {
		params.EventType = "normal"
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rnum := r.Intn(100)

		if channel.DeductionPer >= rnum {
			log.Println("handleApi. channel deduct: test", channel.DeductionPer, rnum)
			params.EventType = "test"
		}

		//ctit
		clkunix, eveunix := clickLogEntity.RequestUnix, params.EventTs/1000000000
		if eveunix-clkunix < int64(30) || (eveunix-clkunix > int64(3*24*60*60)) {
			log.Println("handleApi. cict: test")
			params.EventType = "test"
		}
	}

	if 0 == strings.Compare("test", params.EventType) || 0 == strings.Compare("repeat", params.EventType) {
		go h.saveToDb(params, clickLogEntity, &models.ChannelCbInfo{})
		return
	}

	//回调处理
	cblink := clickLogEntity.CallBack
	if 0 == len(cblink) || -1 == strings.Index(cblink, "http") {
		cblink = channel.Callback
	}

	if 0 == len(cblink) {
		go h.saveToDb(params, clickLogEntity, &models.ChannelCbInfo{})
		return
	}

	cburl := macro.ReplacedClickTLMacroAndFunc(cblink, params, clickLogEntity)
	hstart := time.Now().Unix()
	resCbCont, resCbCode, err := h.ReplacedClickTLSync(cburl)
	errCont := ""
	//落日志clicktoevent eventtoclick stats
	if nil != err {
		errCont = err.Error()
		log.Println(" Callback error spend.", cburl, time.Now().Unix()-hstart, err)

	}

	go h.saveToDb(params, clickLogEntity, &models.ChannelCbInfo{
		CbUrl:  cburl,
		CbCode: resCbCode,
		CbCnt:  resCbCont,
		CbErr:  errCont,
	})
}

func (h *ConvHandler) handleCustomer(params *models.EventParams, eventlogs []models.EventLog, clklog *models.ClickLogEntity) {
	offer, err := h.redisDao.GetCusOffer(params.OfferId)
	//judget event type
	params.EventType = ""
	if 0 < len(eventlogs) {
		log.Println("customer event type repeat because of eventlogs", len(eventlogs), params.ClickId, params.OfferId)
		params.EventType = "repeat"
	}

	if 0 == offer.CampaignId && 0 == len(params.EventType) {
		log.Println("customer event type test because of offer not found", params.ClickId, params.OfferId)
		params.EventType = "test"
	}

	// offcvs, errLog := h.InfoIns.GetOfferStats(clklog.CusOffer.OfferId)

	camcvs, errLog := h.redisDao.GetCampaignDayStats(false, params.OfferId)
	if nil != errLog {
		//错误日志处理
	}

	if camcvs > 0 && 0 == len(params.EventType) && camcvs >= offer.Cap {
		log.Println("customer event type test because of campaign cap", camcvs, params.ClickId, params.OfferId)
		params.EventType = "test"
	}

	if 0 == len(params.EventType) {
		params.EventType = "normal"
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rnum := r.Intn(100)

		deductrate := offer.DeductRate
		if 0 != clklog.CusOffer.AfOfferId && clklog.CusOffer.AfOfferId != clklog.CusOffer.OfferId {
			for _, sw := range offer.Switch {
				if clklog.CusOffer.AfOfferId == sw.TarOfferId {
					deductrate = sw.DeductRate
				}
			}
		}

		if deductrate >= rnum {
			log.Println("customer event type test because of deductrate", deductrate, rnum, params.ClickId, params.OfferId)
			params.EventType = "test"
		}

		//ctit
		clkunix, eveunix := clklog.RequestUnix, params.EventTs/1000000000
		if eveunix-clkunix < int64(30) || (eveunix-clkunix > int64(3*24*60*60)) {
			log.Println("customer event type test because of ctit", eveunix-clkunix, params.ClickId, params.OfferId)
			params.EventType = "test"
		}
	}

	if 0 != strings.Compare("fopen", params.EventName) {
		log.Println("customer event type test because of eventname", params.EventName, params.ClickId, params.OfferId)
		params.EventType = "repeat"
	}

	if 0 == strings.Compare("test", params.EventType) || 0 == strings.Compare("repeat", params.EventType) {
		go h.saveToDb(params, clklog, &models.ChannelCbInfo{})

		return
	}

	//回调处理
	cblink := clklog.CallBack
	// if 0 == len(cblink) || -1 == strings.Index(cblink, "http") {
	// 	cblink = channel.Callback
	// }

	if 0 == len(cblink) {
		go h.saveToDb(params, clklog, &models.ChannelCbInfo{})
		return
	}

	cburl := macro.ReplacedClickTLMacroAndFunc(cblink, params, clklog)
	hstart := time.Now().Unix()
	resCbCont, resCbCode, erre := h.ReplacedClickTLSync(cburl)
	errCont := ""
	//落日志clicktoevent eventtoclick stats
	if nil != erre {
		errCont = erre.Error()
		log.Println(" Callback error spend.", cburl, time.Now().Unix()-hstart, err)

	}

	go h.saveToDb(params, clklog, &models.ChannelCbInfo{
		CbUrl:  cburl,
		CbCode: resCbCode,
		CbCnt:  resCbCont,
		CbErr:  errCont,
	})
}

func (h *ConvHandler) saveToDb(params *models.EventParams, entity *models.ClickLogEntity, info *models.ChannelCbInfo) {

}

func (h *ConvHandler) ReplacedClickTLSync(url string) (string, int, error) {
	return "", 0, nil
}
