package internal

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/models"
	"destroyer-monitor/services/macro"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 转化业务处理
type ApiConvHandler struct {
	ConvMongoDao     *dao.ConvMongoDao
	ClickMongoDaoArr []dao.ClickMongoDao
	RedisDao         *dao.RedisDao
	ClickhouseDao    *dao.ClickhouseDao
}

func (h *ApiConvHandler) Handle(params *models.EventParams, convLogEntity []models.EventLog, clickLogEntity *models.ClickLogEntity) {
	channel, _ := h.RedisDao.GetChannel(strconv.Itoa(clickLogEntity.ChannelId))
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

	offcvs, errLog := h.RedisDao.GetCampaignDayStats(true, params.OfferId)
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

func (h *ApiConvHandler) saveToDb(params *models.EventParams, entity *models.ClickLogEntity, info *models.ChannelCbInfo) {

}

func (h *ApiConvHandler) ReplacedClickTLSync(url string) (string, int, error) {
	return "", 0, nil
}
