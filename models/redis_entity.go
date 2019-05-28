package models

/**
 * 前端通过Redis pub/sub推送过来的数据实体
 */
type ApiOffer struct {
	SalesId      int    `json:"sales_id"`
	SalesName    string `json:"sales_name"`
	OperatorId   int    `json:"operator_id"`
	OperatorName string `json:"operator_name"`

	UniqueId     string   `json:"unique_id"`    //网盟平台订单id
	OfferId      string   `json:"offer_id"`     //网盟平台订单id
	Name         string   `json:"name"`         //订单名称
	PackageName  string   `json:"package_name"` //包名
	PreviewUrl   string   `json:"preview_url"`  //订单的预览链接
	IsIncent     int      `json:"incent"`       //是否支持激励流量
	MinOsv       string   `json:"min_osv"`      //最小系统版本
	InPrice      float64  `json:"price"`
	InCurrency   string   `json:"in_currency"`
	DayCap       int      `json:"day_cap"`
	Countries    []string `json:"countries"`
	Platform     string   `json:"platform"`
	NetTrackLink string   `json:"net_track_link"` //网盟点击链接
	TrackLink    string   `json:"track_link"`     //本平台点击链接
	NetunionId   int      `json:"netunion_id"`    //网盟id
	NetunionName string   `json:"netunion_name"`  //网盟名称
	Updated      int64    `json:"push_ns"`        //更新时间
}

type CusSwitchInfo struct {
	DeductRate int    `json:"deductRate"`
	Priority   int    `json:"priority"`
	Rate       int    `json:"rate"`
	Rule       string `json:"rule"`

	TarSalesId        int     `json:"tarSalesId"`
	TarSalesName      string  `json:"tarSalesName"`
	TarOperatorId     int     `json:"tarOperatorId"`
	TarOperatorName   string  `json:"tarOperatorName"`
	TarAdvertiserId   int     `json:"tarAdvertiserId"`
	TarAdvertiserName string  `json:"tarAdvertiserName"`
	TarCreativeId     int     `json:"tarCreativeId"`
	TarCreativeName   string  `json:"tarCreativeName"`
	TarOfferId        int     `json:"tarOfferId"`
	TarOfferName      string  `json:"tarOfferName"`
	TarInPrice        float64 `json:"tarInPrice"`

	TarPackageName   string `json:"tarPackageName"`
	Os               string `json:"os"`
	CheckIdfa        bool   `json:"checkIdfa"`
	CheckGoogleAdvId bool   `json:"checkGoogleAdvId"`
	CheckAndroidId   bool   `json:"checkAndroidId"`

	TarAdvClikUrl   string `json:"tarAdvClickUrl"`
	TarAdvImpUrl    string `json:"tarAdvImpUrl"`
	TarAdvSyncFunc  string `json:"tarAdvSyncFunc"`
	TarAdvSyncModel string `json:"tarAdvSyncModel"`
}

type CustomerOffer struct {
	SalesId        int    `json:"salesId"`
	SalesName      string `json:"salesName"`
	OperatorId     int    `json:"operatorId"`
	OperatorName   string `json:"operatorName"`
	MediumId       int    `json:"mediumId"`
	MediumName     string `json:"mediumName"`
	AdvertiserId   int    `json:"advertiserId"`
	AdvertiserName string `json:"advertiserName"`
	CreativeId     int    `json:"creativeId"`
	CreativeName   string `json:"creativeName"`

	OfferId      int    `json:"offerId"`
	OfferName    string `json:"offerName"`
	CampaignId   int    `json:"id"`
	CampaignName string `json:"campaignName"`

	ChannelId   int    `json:"channelId"`
	ChannelName string `json:"channelName"`
	ChannelCb   string `json:"channelCallback"`

	Cap      int     `json:"cap"`
	OfferCap int     `json:"offerCap"`
	InPrice  float64 `json:"inPrice"`
	OutPrice float64 `json:"outPrice"`

	PackageName      string `json:"packageName"`
	Os               string `json:"os"`
	CheckIdfa        bool   `json:"checkIdfa"`
	CheckGoogleAdvId bool   `json:"checkGoogleAdvId"`
	CheckAndroidId   bool   `json:"checkAndroidId"`

	DeductRate   int    `json:"deductRate"`
	ClickUrl     string `json:"clkUrl"`
	AdvClikUrl   string `json:"advClickUrl"`
	AdvImpUrl    string `json:"advImpUrl"`
	AdvSyncFunc  string `json:"advSyncFunc"`
	AdvSyncModel string `json:"advSyncModel"`

	CampaignSyncMode string `json:"campaignSyncModel"`

	DisSubchannels []string `json:"disableSubChannels"`

	Switch []CusSwitchInfo `json:"offerSwitch"`
}

type Channel struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Callback      string `json:"callback_url"`
	PriceDiscount int    `json:"price_discount"`
	DeductionPer  int    `json:"deduction_percentage"`
	Media         int    `json:"media"`
	MediaName     string `json:"media_name"`
}
