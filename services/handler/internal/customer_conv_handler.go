package internal

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/models"
	"destroyer-monitor/services/macro"
	"log"
	"math/rand"
	"strings"
	"time"
)

// 转化业务处理
type CustomerConvHandler struct {
	ConvMongoDao     *dao.ConvMongoDao
	ClickMongoDaoArr []dao.ClickMongoDao
	RedisDao         *dao.RedisDao
	ClickhouseDao    *dao.ClickhouseDao
}

func (h *CustomerConvHandler) Handle(params *models.EventParams, convLogEntity []models.EventLog, clickLogEntity *models.ClickLogEntity) {
	offer, err := h.RedisDao.GetCusOffer(params.OfferId)
	//judget event type
	params.EventType = ""
	if 0 < len(convLogEntity) {
		log.Println("customer event type repeat because of convLogEntity", len(convLogEntity), params.ClickId, params.OfferId)
		params.EventType = "repeat"
	}

	if 0 == offer.CampaignId && 0 == len(params.EventType) {
		log.Println("customer event type test because of offer not found", params.ClickId, params.OfferId)
		params.EventType = "test"
	}

	// offcvs, errLog := h.InfoIns.GetOfferStats(clickLogEntity.CusOffer.OfferId)

	camcvs, errLog := h.RedisDao.GetCampaignDayStats(false, params.OfferId)
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
		if 0 != clickLogEntity.CusOffer.AfOfferId && clickLogEntity.CusOffer.AfOfferId != clickLogEntity.CusOffer.OfferId {
			for _, sw := range offer.Switch {
				if clickLogEntity.CusOffer.AfOfferId == sw.TarOfferId {
					deductrate = sw.DeductRate
				}
			}
		}

		if deductrate >= rnum {
			log.Println("customer event type test because of deductrate", deductrate, rnum, params.ClickId, params.OfferId)
			params.EventType = "test"
		}

		//ctit
		clkunix, eveunix := clickLogEntity.RequestUnix, params.EventTs/1000000000
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
		go h.saveToDb(params, clickLogEntity, &models.ChannelCbInfo{})

		return
	}

	//回调处理
	cblink := clickLogEntity.CallBack
	// if 0 == len(cblink) || -1 == strings.Index(cblink, "http") {
	// 	cblink = channel.Callback
	// }

	if 0 == len(cblink) {
		go h.saveToDb(params, clickLogEntity, &models.ChannelCbInfo{})
		return
	}

	cburl := macro.ReplacedClickTLMacroAndFunc(cblink, params, clickLogEntity)
	hstart := time.Now().Unix()
	resCbCont, resCbCode, erre := h.ReplacedClickTLSync(cburl)
	errCont := ""
	//落日志clicktoevent eventtoclick stats
	if nil != erre {
		errCont = erre.Error()
		log.Println(" Callback error spend.", cburl, time.Now().Unix()-hstart, err)

	}

	go h.saveToDb(params, clickLogEntity, &models.ChannelCbInfo{
		CbUrl:  cburl,
		CbCode: resCbCode,
		CbCnt:  resCbCont,
		CbErr:  errCont,
	})
}

func (h *CustomerConvHandler) saveToDb(params *models.EventParams, entity *models.ClickLogEntity, info *models.ChannelCbInfo) {

}

func (h *CustomerConvHandler) ReplacedClickTLSync(url string) (string, int, error) {
	return "", 0, nil
}
