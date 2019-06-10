package internal

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/lib/queue"
	"destroyer-monitor/models"
	"destroyer-monitor/utils"
	"encoding/json"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CustomerClickHandler struct {
	RedisDao *dao.RedisDao
	Queue    queue.Queue
	HttpCli  *utils.HttpConPool
}

func (c *CustomerClickHandler) Handle(params *models.ClickParams) *models.Response {
	offer, err := c.RedisDao.GetCusOffer(params.OfferId)
	if nil != err {
		//日志
		return &models.Response{
			Code:    200,
			Content: "parameter error [offer_id]",
		}
	}

	daycvs, err := c.RedisDao.GetDayCap(false, strconv.Itoa(offer.OfferId), params.OfferId)
	if nil != err && err.ErrType != models.RedisQueryNil {
		//错误处理
		// go c.CatchLog.SysErrorLog(params)
	}

	clickId := strings.Join([]string{
		params.ChannelId,
		params.OfferId,
		strconv.Itoa(offer.OfferId),
		strconv.FormatInt(params.ClickTs, 10)}, "^")
	params.ClickId = utils.GenUniqueIdExt(params.OfferId, params.ClickTs, utils.Md5(clickId))

	switchinfo := c.judgeCampaignSwitch(offer, daycvs)

	fileter, fcause := c.filterCampaign(params, offer, switchinfo)
	params.InvalidCause = fcause
	if fileter {
		//日志
		go c.CusLogFilter(params, offer, switchinfo)

		return &models.Response{
			Code:    200,
			Content: fcause,
		}
	}

	return c.postProcess(params, offer, switchinfo)
}

func (c *CustomerClickHandler) postProcess(params *models.ClickParams, offer *models.CustomerOffer, switchinfo *models.CusSwitchInfo) *models.Response {

	//切换
	if nil != switchinfo {
		params.OffSyncType = switchinfo.TarAdvSyncModel
		if 0 == strings.Compare("client", switchinfo.TarAdvSyncModel) {

			srclink := switchinfo.TarAdvClikUrl

			relink := strings.Replace(srclink, "__CLICKID__", params.ClickId, -1)
			relink = strings.Replace(relink, "__IDFA__", params.Idfa, -1)
			relink = strings.Replace(relink, "__GOOGLEADVID__", params.GoogleAdvId, -1)
			relink = strings.Replace(relink, "__ANDROIDID__", params.AndroidId, -1)
			relink = strings.Replace(relink, "__SUBCHANNEL__", params.SubChannel, -1)
			relink = strings.Replace(relink, "__CHANNEL__", params.ChannelId, -1)
			relink = strings.Replace(relink, "__IP__", params.Ip, -1)
			relink = strings.Replace(relink, "__UA__", url.QueryEscape(params.Ua), -1)
			relink = strings.Replace(relink, "__S1__", params.S1, -1)
			relink = strings.Replace(relink, "__S2__", params.S2, -1)
			relink = strings.Replace(relink, "__S3__", params.S3, -1)
			relink = strings.Replace(relink, "__S4__", params.S4, -1)
			relink = strings.Replace(relink, "__S5__", params.S5, -1)

			params.ReplaceLink = relink
			go c.CusOfferSuccess(params, offer, switchinfo)
			//日志
			//
			return &models.Response{
				Code:    302,
				Content: relink,
			}
		}

		if 0 == strings.Compare("server", switchinfo.TarAdvSyncModel) {
			go c.syncAdvLinkProcess(params, offer, switchinfo)

			return &models.Response{
				Code:    200,
				Content: "success",
			}
		}
	}

	//正常

	params.OffSyncType = offer.AdvSyncModel
	if 0 == strings.Compare("client", offer.AdvSyncModel) {
		srclink := offer.AdvClikUrl
		relink := strings.Replace(srclink, "__CLICKID__", params.ClickId, -1)
		relink = strings.Replace(relink, "__IDFA__", params.Idfa, -1)
		relink = strings.Replace(relink, "__GOOGLEADVID__", params.GoogleAdvId, -1)
		relink = strings.Replace(relink, "__ANDROIDID__", params.AndroidId, -1)
		relink = strings.Replace(relink, "__SUBCHANNEL__", params.SubChannel, -1)
		relink = strings.Replace(relink, "__CHANNEL__", params.ChannelId, -1)
		relink = strings.Replace(relink, "__IP__", params.Ip, -1)
		relink = strings.Replace(relink, "__UA__", url.QueryEscape(params.Ua), -1)
		relink = strings.Replace(relink, "__S1__", params.S1, -1)
		relink = strings.Replace(relink, "__S2__", params.S2, -1)
		relink = strings.Replace(relink, "__S3__", params.S3, -1)
		relink = strings.Replace(relink, "__S4__", params.S4, -1)
		relink = strings.Replace(relink, "__S5__", params.S5, -1)

		params.ReplaceLink = relink
		go c.CusOfferSuccess(params, offer, switchinfo)

		return &models.Response{
			Code:    302,
			Content: relink,
		}
	}

	if 0 == strings.Compare("server", offer.AdvSyncModel) {
		go c.syncAdvLinkProcess(params, offer, switchinfo)

		return &models.Response{
			Code:    200,
			Content: "success",
		}
	}

	return &models.Response{
		Code:    404,
		Content: "NotFound",
	}
}

func (c *CustomerClickHandler) syncAdvLinkProcess(params *models.ClickParams, offer *models.CustomerOffer, switchinfo *models.CusSwitchInfo) {
	link := offer.AdvClikUrl

	if nil != switchinfo {
		link = switchinfo.TarAdvClikUrl
	}

	relink := strings.Replace(link, "__CLICKID__", params.ClickId, -1)
	relink = strings.Replace(relink, "__IDFA__", params.Idfa, -1)
	relink = strings.Replace(relink, "__GOOGLEADVID__", params.GoogleAdvId, -1)
	relink = strings.Replace(relink, "__ANDROIDID__", params.AndroidId, -1)
	relink = strings.Replace(relink, "__SUBCHANNEL__", params.SubChannel, -1)
	relink = strings.Replace(relink, "__CHANNEL__", params.ChannelId, -1)

	cnt, code, err := c.HttpCli.ReqGet(relink)

	params.ReplaceLink = relink
	params.SyncCode = code
	if 512 > len(cnt) {
		params.SyncCnt = cnt
	} else {
		params.SyncCnt = string([]rune(cnt)[:512])
	}
	if nil != err {
		params.SyncErr = err.Error()
	}
	//日志
	go c.CusOfferSuccess(params, offer, switchinfo)
}

func (c *CustomerClickHandler) filterCampaign(params *models.ClickParams, offer *models.CustomerOffer, switchinfo *models.CusSwitchInfo) (bool, string) {

	if nil != switchinfo {
		//parmas
		if switchinfo.CheckIdfa && 0 == len(params.Idfa) {
			return true, "error matched [idfa]"
		}
		if switchinfo.CheckAndroidId && 0 == len(params.AndroidId) {
			return true, "error matched [androidid]"
		}
		if switchinfo.CheckGoogleAdvId && 0 == len(params.GoogleAdvId) {
			return true, "error matched [google advertiser id]"
		}

		return false, ""
	}

	if 0 < len(params.SubChannel) {
		for _, subc := range offer.DisSubchannels {
			if 0 == strings.Compare(params.SubChannel, subc) {
				return true, "error matched [disable subchannel]"
			}
		}
	}

	//TODO:os 不过滤
	// if (0 == strings.Compare(offer.CampaignSyncMode, "client")) && (0 != strings.Compare(strings.ToUpper(params.RequestOs), strings.ToUpper(offer.Os))) {
	// 	return true, "error matched [os]"
	// }

	//parmas
	if offer.CheckIdfa && 0 == len(params.Idfa) {
		return true, "error matched [idfa]"
	}
	if offer.CheckAndroidId && 0 == len(params.AndroidId) {
		return true, "error matched [androidid]"
	}
	if offer.CheckGoogleAdvId && 0 == len(params.GoogleAdvId) {
		return true, "error matched [google advertiser id]"
	}

	return false, ""
}

func (c *CustomerClickHandler) judgeCampaignSwitch(offer *models.CustomerOffer, daycvs map[string]int) *models.CusSwitchInfo {
	//切换优先级排序
	for i := 1; i < len(offer.Switch); i++ {
		j := i - 1
		key := offer.Switch[i]
		for {
			if j >= 0 && offer.Switch[j].Priority > key.Priority {
				offer.Switch[j+1] = offer.Switch[j]
				j = j - 1
			} else {
				break
			}
		}
		offer.Switch[j+1] = key
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnum := r.Intn(100)

	for _, swinfo := range offer.Switch {
		if 0 == strings.Compare(swinfo.Rule, "random") && swinfo.Rate > rnum {
			return &swinfo
		}

		if 0 == strings.Compare(swinfo.Rule, "ccap") && daycvs["campaign"] > int(offer.Cap*swinfo.Rate/100) {
			return &swinfo
		}

		if 0 == strings.Compare(swinfo.Rule, "ocap") && daycvs["offer"] > int(offer.OfferCap*swinfo.Rate/100) {
			return &swinfo
		}
	}

	return nil
}

func (l *CustomerClickHandler) CusLogFilter(params *models.ClickParams, offer *models.CustomerOffer, swoffer *models.CusSwitchInfo) {
	logs := models.ClickLog{
		LogType:   "click",
		OfferType: "customer",

		ClickId: params.ClickId,

		ChannelId:    offer.ChannelId,
		ChannelName:  offer.ChannelName,
		SubChannel:   params.SubChannel,
		ChannelClkId: params.ChannelClkId,

		DestroyerId:   params.OfferId,
		DestroyerName: offer.CampaignName,

		IsValid:      false,
		InvalidCause: params.InvalidCause,

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

		OffSyncType:     offer.CampaignSyncMode,
		DirectUrl:       params.ReplaceLink,
		SyncAdvLink:     params.ReplaceLink,
		SyncAdvSycnCode: params.SyncCode,
		SyncAdvSycnText: params.SyncCnt,
	}

	cusoffer := models.CusOfferStr{
		MediumId:     offer.MediumId,
		MediumName:   offer.MediumName,
		OperatorId:   offer.OperatorId,
		OperatorName: offer.OperatorName,
		SalesId:      offer.SalesId,
		SalesName:    offer.SalesName,
		AdvId:        offer.AdvertiserId,
		AdvName:      offer.AdvertiserName,
		OfferId:      offer.OfferId,
		OfferName:    offer.OfferName,
		PackageName:  offer.PackageName,
		CreativeId:   offer.CreativeId,
		CreativeName: offer.CreativeName,
		Platform:     strings.ToUpper(offer.Os),
		InPrice:      offer.InPrice,
		InCurrency:   "USD",
		OutPrice:     offer.OutPrice,
		OutCurrency:  "USD",
	}

	if nil != swoffer {
		logs.IsSwitch = true

		cusoffer.AfOperatorId = swoffer.TarOperatorId
		cusoffer.AfOperatorName = swoffer.TarOperatorName
		cusoffer.AfSalesId = swoffer.TarSalesId
		cusoffer.AfSalesName = swoffer.TarSalesName
		cusoffer.AfAdvId = swoffer.TarAdvertiserId
		cusoffer.AfAdvName = swoffer.TarAdvertiserName
		cusoffer.AfOfferId = swoffer.TarOfferId
		cusoffer.AfofferName = swoffer.TarOfferName
		cusoffer.AfCreativeId = swoffer.TarCreativeId
		cusoffer.AfCreativeName = swoffer.TarCreativeName
		cusoffer.AfPlatform = strings.ToUpper(swoffer.Os)
		cusoffer.AfInPrice = swoffer.TarInPrice
		cusoffer.AfInCurrency = "USD"
	}
	logs.CusOffer = cusoffer

	if 0 == len(logs.CallBack) && 0 < len(offer.ChannelCb) {
		logs.CallBack = offer.ChannelCb
	}

	bytes, err := json.Marshal(logs)
	if nil != err {
		log.Println("SysErrorLog marshal data error.", err.Error())
		return
	}

	_ = l.Queue.Enqueue(logs.ClickId, bytes)
}

func (l *CustomerClickHandler) CusOfferSuccess(params *models.ClickParams, offer *models.CustomerOffer, swoffer *models.CusSwitchInfo) {
	logs := models.ClickLog{
		LogType:   "click",
		OfferType: "customer",

		ClickId: params.ClickId,

		ChannelId:    offer.ChannelId,
		ChannelName:  offer.ChannelName,
		SubChannel:   params.SubChannel,
		ChannelClkId: params.ChannelClkId,

		DestroyerId:   params.OfferId,
		DestroyerName: offer.CampaignName,

		IsValid:      true,
		InvalidCause: params.InvalidCause,

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

		OffSyncType:     offer.CampaignSyncMode,
		DirectUrl:       params.ReplaceLink,
		SyncAdvLink:     params.ReplaceLink,
		SyncAdvSycnCode: params.SyncCode,
		SyncAdvSycnText: params.SyncCnt,
	}

	cusoffer := models.CusOfferStr{
		MediumId:     offer.MediumId,
		MediumName:   offer.MediumName,
		OperatorId:   offer.OperatorId,
		OperatorName: offer.OperatorName,
		SalesId:      offer.SalesId,
		SalesName:    offer.SalesName,
		AdvId:        offer.AdvertiserId,
		AdvName:      offer.AdvertiserName,
		OfferId:      offer.OfferId,
		OfferName:    offer.OfferName,
		PackageName:  offer.PackageName,
		Platform:     strings.ToUpper(offer.Os),
		InPrice:      offer.InPrice,
		InCurrency:   "USD",
		OutPrice:     offer.OutPrice,
		OutCurrency:  "USD",
		CreativeId:   offer.CreativeId,
		CreativeName: offer.CreativeName,
	}

	if nil != swoffer {
		logs.IsSwitch = true

		cusoffer.AfOperatorId = swoffer.TarOperatorId
		cusoffer.AfOperatorName = swoffer.TarOperatorName
		cusoffer.AfSalesId = swoffer.TarSalesId
		cusoffer.AfSalesName = swoffer.TarSalesName
		cusoffer.AfAdvId = swoffer.TarAdvertiserId
		cusoffer.AfAdvName = swoffer.TarAdvertiserName
		cusoffer.AfOfferId = swoffer.TarOfferId
		cusoffer.AfofferName = swoffer.TarOfferName
		cusoffer.AfCreativeId = swoffer.TarCreativeId
		cusoffer.AfCreativeName = swoffer.TarCreativeName
		cusoffer.AfPlatform = strings.ToUpper(swoffer.Os)
		cusoffer.AfInPrice = swoffer.TarInPrice
		cusoffer.AfInCurrency = "USD"
	}
	logs.CusOffer = cusoffer
	if 0 == len(logs.CallBack) && 0 < len(offer.ChannelCb) {
		logs.CallBack = offer.ChannelCb
	}

	bytes, err := json.Marshal(logs)
	if nil != err {
		log.Println("SysErrorLog marshal data error.", err.Error())
		return
	}

	_ = l.Queue.Enqueue(logs.ClickId, bytes)
}
