package internal

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/models"
	"destroyer-monitor/services"
	"destroyer-monitor/utils"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type ApiClickHandler struct {
	RedisDao *dao.RedisDao
	Queue    queue.Queue
}

func (h *ApiClickHandler) Handle(params *models.ClickParams) *models.Response {
	netlenstr := string([]rune(params.OfferId)[:1])
	netlen, err := strconv.ParseInt(netlenstr, 10, 32)
	if nil != err {
		return &models.Response{
			Code:    200,
			Content: "parameter error [offer_id]",
		}
	}
	netunionid := string([]rune(params.OfferId)[1 : 1+netlen])
	_, err = strconv.ParseInt(netunionid, 10, 32)
	if nil != err {
		return &models.Response{
			Code:    200,
			Content: "parameter error [offer_id]",
		}
	}

	channel, errt := h.RedisDao.GetChannel(params.ChannelId)
	if nil != errt {
		zap.Get().Error("get channel error.", params.ChannelId, errt)
		return &models.Response{
			Code:    200,
			Content: "parameter error [channel not exists]",
		}
	}

	offer, errt := h.RedisDao.GetApiOffer(netunionid, params.OfferId)
	if nil != errt {
		//订单获取失败，走切换流程
		clickId := strings.Join([]string{params.ChannelId, params.OfferId, strconv.FormatInt(params.ClickTs, 10)}, "^")
		params.ClickId = utils.GenUniqueIdExt(params.OfferId, params.ClickTs, utils.Md5(clickId))
		return h.switchToSmartlink(params, nil, channel, "offer not exists")
	}

	blackListInfo, errt := h.RedisDao.GetBlacklist(params.OfferId)
	if nil != errt && errt.ErrType != models.RedisQueryNil {
		log.Println("blackListInfo info. error.", blackListInfo, errt)
		zap.Get().Error("query black list error, params.OfferId=", params.OfferId, err)
	}

	daycvs, errt := h.RedisDao.GetDayCap(true, "", params.OfferId)
	if nil != errt && errt.ErrType != models.RedisQueryNil {
		log.Println("daycvs info. error.", daycvs, errt)
	}

	return h.process(params, offer, channel, blackListInfo, daycvs)
}

func (h *ApiClickHandler) process(params *models.ClickParams, offer *models.ApiOffer,
	channel *models.Channel, blackListInfos map[string]string, daycvs map[string]int) *models.Response {
	clickId := strings.Join([]string{params.ChannelId, params.OfferId, strconv.Itoa(offer.NetunionId), strconv.FormatInt(params.ClickTs, 10)}, "^")
	params.ClickId = utils.GenUniqueIdExt(params.OfferId, params.ClickTs, utils.Md5(clickId))

	//达到cap后，进入切换流程
	if offer.DayCap > 0 && offer.DayCap-daycvs["campaign"] <= 0 {
		return h.switchToSmartlink(params, offer, channel, "up to offer cap")
	}

	//检查黑名单
	if ty, ok := blackListInfos["type"]; ok {

		//全局黑名单，进入切换流程
		if 0 == strings.Compare("0", ty) {
			return h.switchToSmartlink(params, offer, channel, "offer disable.")
		}

		//订单对渠道级别黑名单，进入切换流程
		if _, ok := blackListInfos[strings.Join([]string{"channel", params.ChannelId}, ":")]; ok {
			return h.switchToSmartlink(params, offer, channel, "offer disable..")
		}

		//订单对子渠道渠道级别黑名单，进入切换流程
		if _, ok := blackListInfos[strings.Join([]string{"channel", params.ChannelId, "subchannel", params.SubChannel}, ":")]; ok {
			return h.switchToSmartlink(params, offer, channel, "offer disable...")
		}
	}

	//宏替换后，记录日志，并返回
	srclink := offer.NetTrackLink
	relink := strings.Replace(srclink, "__CLICKID__", params.ClickId, -1)
	relink = strings.Replace(relink, "__IDFA__", params.Idfa, -1)
	relink = strings.Replace(relink, "__GOOGLEADVID__", params.GoogleAdvId, -1)
	relink = strings.Replace(relink, "__ANDROIDID__", params.AndroidId, -1)
	relink = strings.Replace(relink, "__SUBCHANNEL__", params.SubChannel, -1)
	relink = strings.Replace(relink, "__CHANNEL__", params.ChannelId, -1)

	params.ReplaceLink = relink

	go h.apiLogSuccess(params, offer, channel)

	return &models.Response{
		Code:    302,
		Content: relink,
	}
}

func (h *ApiClickHandler) switchToSmartlink(params *models.ClickParams, offer *models.ApiOffer, channel *models.Channel, cause string) *models.Response {
	if 0 < len(services.SmarklinkOffers) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rindex := r.Intn(len(services.SmarklinkOffers))

		campaignid := services.SmarklinkOffers[rindex]
		netlenstr := string([]rune(campaignid)[:1])
		netlen, err := strconv.ParseInt(netlenstr, 10, 32)
		if nil != err {
			// go h.CatchLog.SysErrorLog(params)
			log.Println("switchToSmartlink. parse netunion id error.", campaignid, err.Error())
			return &models.Response{
				Code:    200,
				Content: "parameter error [.offer_id]",
			}
		}
		netunionid := string([]rune(campaignid)[1 : 1+netlen])

		smoffer, errt := h.RedisDao.GetApiOffer(netunionid, campaignid)
		if nil != errt {
			log.Println("switchToSmartlink. get smartlink offer error.", netunionid, campaignid, errt)
			return &models.Response{
				Code:    200,
				Content: "parameter error [.offer_id.]",
			}
		}

		srclink := smoffer.NetTrackLink
		relink := strings.Replace(srclink, "__CLICKID__", params.ClickId, -1)
		relink = strings.Replace(relink, "__IDFA__", params.Idfa, -1)
		relink = strings.Replace(relink, "__GOOGLEADVID__", params.GoogleAdvId, -1)
		relink = strings.Replace(relink, "__ANDROIDID__", params.AndroidId, -1)
		relink = strings.Replace(relink, "__SUBCHANNEL__", params.SubChannel, -1)
		relink = strings.Replace(relink, "__CHANNEL__", params.ChannelId, -1)

		params.ReplaceLink = relink

		go h.apiLogSwitchSuccess(params, offer, smoffer, channel)

		return &models.Response{
			Code:    302,
			Content: relink,
		}
	}

	return &models.Response{
		Code:    200,
		Content: strings.Join([]string{"error", cause}, ":"),
	}
}

func (l *ApiClickHandler) apiLogSuccess(params *models.ClickParams, offer *models.ApiOffer, channel *models.Channel) {

	logs := models.ClickLog{
		LogType:   "click",
		OfferType: "api",

		ClickId: params.ClickId,

		ChannelId:    channel.Id,
		ChannelName:  channel.Name,
		SubChannel:   params.SubChannel,
		ChannelClkId: params.ChannelClkId,

		IsSwitch: false,

		DestroyerId:   params.OfferId,
		DestroyerName: offer.Name,

		ApiOffer: models.ApiOfferStr{
			MediumId:       channel.Media,
			MediumName:     channel.MediaName,
			OperatorId:     offer.OperatorId,
			OperatorName:   offer.OperatorName,
			SalesId:        offer.SalesId,
			SalesName:      offer.SalesName,
			NetunionId:     offer.NetunionId,
			NetunionName:   offer.NetunionName,
			OffferUniqueId: params.OfferId,
			OfferId:        offer.OfferId,
			OfferName:      offer.Name,
			PackageName:    offer.PackageName,
			Platform:       strings.ToUpper(offer.Platform),
			InPrice:        offer.InPrice,
			InCurrency:     offer.InCurrency,
			OutPrice:       offer.InPrice * float64(channel.PriceDiscount) / float64(100),
			OutCurrency:    offer.InCurrency,
			Country:        strings.Join(offer.Countries, ","),
		},

		ApiOfferCap: offer.DayCap,

		IsValid:      true,
		InvalidCause: "",

		ReqestUrl:      params.RequestUrl,
		RequestUnix:    time.Now().Unix(),
		RequestTime:    time.Now().Format("2006-01-02 15:04:05"),
		RequestIp:      params.RequestIp,
		RequestCountry: params.RequestCountry,
		RequestUa:      params.RequestUa,
		RequestOs:      strings.ToUpper(params.RequestOs),

		Ip:          params.Ip,
		Ua:          params.Ua,
		Idfa:        params.Idfa,
		AndroidId:   params.AndroidId,
		GoogleAdvId: params.GoogleAdvId,
		Imei:        params.Imei,
		Os:          strings.ToUpper(params.RequestOs),
		CallBack:    params.Callback,
		S1:          params.S1,
		S2:          params.S2,
		S3:          params.S3,
		S4:          params.S4,
		S5:          params.S5,

		OffSyncType:     "client",
		DirectUrl:       params.ReplaceLink,
		SyncAdvLink:     params.ReplaceLink,
		SyncAdvSycnCode: params.SyncCode,
		SyncAdvSycnText: params.SyncCnt,
	}

	bytes, err := json.Marshal(logs)
	if nil != err {
		log.Println("SysErrorLog marshal data error.", err.Error())
		return
	}

	_ = l.Queue.Enqueue(logs.ClickId, bytes)
}

func (l *ApiClickHandler) apiLogSwitchSuccess(params *models.ClickParams, offer *models.ApiOffer, swoffer *models.ApiOffer, channel *models.Channel) {
	logs := models.ClickLog{
		LogType:   "click",
		OfferType: "api",

		ClickId: params.ClickId,

		// ChannelId:    channel.Id,
		// ChannelName:  channel.Name,
		SubChannel:   params.SubChannel,
		ChannelClkId: params.ChannelClkId,

		IsSwitch: true,

		DestroyerId: params.OfferId,
		// DestroyerName: offer.Name,

		ApiOffer: models.ApiOfferStr{
			AfOperatorId:     swoffer.OperatorId,
			AfOperatorName:   swoffer.OperatorName,
			AfSalesId:        swoffer.SalesId,
			AfSalesName:      swoffer.SalesName,
			AfNetunionId:     swoffer.NetunionId,
			AfNetunionName:   swoffer.NetunionName,
			AfOffferUniqueId: swoffer.UniqueId,
			AfOfferId:        swoffer.OfferId,
			AfOfferName:      swoffer.Name,
			AfPackageName:    swoffer.PackageName,
			AfPlatform:       strings.ToUpper(swoffer.Platform),
			AfInPrice:        swoffer.InPrice,
			AfInCurrency:     swoffer.InCurrency,
			AfOutCurrency:    swoffer.InCurrency,
			AfCountry:        strings.Join(swoffer.Countries, ","),
		},

		IsValid:      true,
		InvalidCause: "",

		ReqestUrl:      params.RequestUrl,
		RequestUnix:    time.Now().Unix(),
		RequestTime:    time.Now().Format("2006-01-02 15:04:05"),
		RequestIp:      params.RequestIp,
		RequestCountry: params.RequestCountry,
		RequestUa:      params.RequestUa,
		RequestOs:      strings.ToUpper(params.RequestOs),

		Ip:          params.Ip,
		Ua:          params.Ua,
		Idfa:        params.Idfa,
		AndroidId:   params.AndroidId,
		GoogleAdvId: params.GoogleAdvId,
		Imei:        params.Imei,
		Os:          strings.ToUpper(params.RequestOs),
		CallBack:    params.Callback,
		S1:          params.S1,
		S2:          params.S2,
		S3:          params.S3,
		S4:          params.S4,
		S5:          params.S5,

		OffSyncType:     "client",
		DirectUrl:       params.ReplaceLink,
		SyncAdvLink:     params.ReplaceLink,
		SyncAdvSycnCode: params.SyncCode,
		SyncAdvSycnText: params.SyncCnt,
	}

	if channel != nil {
		logs.ApiOffer.MediumId = channel.Media
		logs.ApiOffer.MediumName = channel.MediaName
		logs.ApiOffer.AfOutPrice = swoffer.InPrice * float64(channel.PriceDiscount) / float64(100)
		logs.ChannelId = channel.Id
		logs.ChannelName = channel.Name
	}

	logs.ApiOffer.OffferUniqueId = params.OfferId
	if offer != nil {
		logs.ApiOffer.OperatorId = offer.OperatorId
		logs.ApiOffer.OperatorName = offer.OperatorName
		logs.ApiOffer.SalesId = offer.SalesId
		logs.ApiOffer.SalesName = offer.SalesName
		logs.ApiOffer.NetunionId = offer.NetunionId
		logs.ApiOffer.NetunionName = offer.NetunionName
		logs.ApiOffer.OfferId = offer.OfferId
		logs.ApiOffer.OfferName = offer.Name
		logs.ApiOffer.PackageName = offer.PackageName
		logs.ApiOffer.Platform = strings.ToUpper(offer.Platform)
		logs.ApiOffer.InPrice = offer.InPrice
		logs.ApiOffer.InCurrency = offer.InCurrency
		logs.ApiOffer.OutPrice = offer.InPrice * float64(channel.PriceDiscount) / float64(100)
		logs.ApiOffer.OutCurrency = offer.InCurrency
		logs.ApiOffer.Country = strings.Join(offer.Countries, ",")
	}

	bytes, err := json.Marshal(logs)
	if nil != err {
		log.Println("SysErrorLog marshal data error.", err.Error())
		return
	}

	_ = l.Queue.Enqueue(logs.ClickId, bytes)
}
