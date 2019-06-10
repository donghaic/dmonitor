package models

/**
 * 对应mongo数据库内的数据实体
 */

type ClickLogEntity struct {
	LogType   string `bson:"log_type"`
	OfferType string `bson:"offer_type"` //"customer"  "api"

	ClickId string `bson:"clickid"`

	ChannelId    int    `bson:"channel_id"`
	ChannelName  string `bson:"channel_name"`
	SubChannel   string `bson:"sub_channel"`
	ChannelClkId string `bson:"channel_clkid"`

	IsSwitch bool `bson:"is_switch"`

	DestroyerId   string `bson:"destroyer_id"` //对外的订单id  直客是campaignid api是订单独立id
	DestroyerName string `bson:"destroyer_name"`

	ApiOffer struct {
		MediumId       int     `bson:"medium_id"`
		MediumName     string  `bson:"medium_name"`
		OperatorId     int     `bson:"operator_id"`
		OperatorName   string  `bson:"operator_name"`
		SalesId        int     `bson:"sales_id"`
		SalesName      string  `bson:"sales_name"`
		NetunionId     int     `bson:"netunion_id"`
		NetunionName   string  `bson:"netunion_name"`
		OffferUniqueId string  `bson:"offer_uniqueid"`
		OfferId        string  `bson:"offer_id"`
		OfferName      string  `bson:"offer_name"`
		PackageName    string  `bson:"package_name"`
		Platform       string  `bson:"platform"`
		Country        string  `bson:"country"`
		InPrice        float64 `bson:"in_price"`
		InCurrency     string  `bson:"in_currency"`
		OutPrice       float64 `bson:"out_price"`
		OutCurrency    string  `bson:"out_currency"`

		AfOperatorId     int     `bson:"afoperator_id"`
		AfOperatorName   string  `bson:"afoperator_name"`
		AfSalesId        int     `bson:"afsales_id"`
		AfSalesName      string  `bson:"afsales_name"`
		AfNetunionId     int     `bson:"afnetunion_id"`
		AfNetunionName   string  `bson:"afnetunion_name"`
		AfOffferUniqueId string  `bson:"afoffer_uniqueid"`
		AfOfferId        string  `bson:"afoffer_id"`
		AfOfferName      string  `bson:"afoffer_name"`
		AfPackageName    string  `bson:"afpackage_name"`
		AfPlatform       string  `bson:"afplatform"`
		AfCountry        string  `bson:"afcountry"`
		AfInPrice        float64 `bson:"afin_price"`
		AfInCurrency     string  `bson:"afin_currency"`
		AfOutPrice       float64 `bson:"afout_price"`
		AfOutCurrency    string  `bson:"afout_currency"`
	} `bson:"api_offer"`
	CusOffer struct {
		MediumId     int     `bson:"medium_id"`
		MediumName   string  `bson:"medium_name"`
		OperatorId   int     `bson:"operator_id"`
		OperatorName string  `bson:"operator_name"`
		SalesId      int     `bson:"sales_id"`
		SalesName    string  `bson:"sales_name"`
		AdvId        int     `bson:"adv_id"`
		AdvName      string  `bson:"adv_name"`
		OfferId      int     `bson:"offer_id"`
		OfferName    string  `bson:"offer_name"`
		CreativeId   int     `bson:"creative_id"`
		CreativeName string  `bson:"creative_name"`
		PackageName  string  `bson:"package_name"`
		Platform     string  `bson:"platform"`
		InPrice      float64 `bson:"in_price"`
		InCurrency   string  `bson:"in_currency"`
		OutPrice     float64 `bson:"out_price"`
		OutCurrency  string  `bson:"out_currency"`

		AfOperatorId   int     `bson:"afoperator_id"`
		AfOperatorName string  `bson:"afoperator_name"`
		AfSalesId      int     `bson:"afsales_id"`
		AfSalesName    string  `bson:"afsales_name"`
		AfAdvId        int     `bson:"afadv_id"`
		AfAdvName      string  `bson:"afadv_name"`
		AfOfferId      int     `bson:"afoffer_id"`
		AfofferName    string  `bson:"afoffer_name"`
		AfCreativeId   int     `bson:"afcreative_id"`
		AfCreativeName string  `bson:"afcreative_name"`
		AfPackageName  string  `bson:"afpackage_name"`
		AfPlatform     string  `bson:"afplatform"`
		AfInPrice      float64 `bson:"afin_price"`
		AfInCurrency   string  `bson:"afout_currency"`
	} `bson:"customer_offer"`

	OfferCap    int `bson:"offer_cap"`
	CampaignCap int `bson:"campaign_cap"`

	IsValid      bool   `bson:"is_valid"`
	InvalidCause string `bson:"invalid_cause"`

	ReqestUrl      string `bson:"request_url"`
	RequestUnix    int64  `bson:"requestunix"`
	RequestTime    string `bson:"request_time"`
	RequestIp      string `bson:"request_ip"`
	RequestCountry string `bson:"request_country"`
	RequestUa      string `bson:"request_ua"`
	RequestOs      string `bson:"request_os"`

	Ip          string `bson:"ip"`
	Ua          string `bson:"ua"`
	Idfa        string `bson:"idfa"`
	AndroidId   string `bson:"androidid"`
	GoogleAdvId string `bson:"google_advid"`
	Imei        string `bson:"imei"`
	Os          string `bson:"os"`
	CallBack    string `bson:"callback"`
	S1          string `bson:"s1"`
	S2          string `bson:"s2"`
	S3          string `bson:"s3"`
	S4          string `bson:"s4"`
	S5          string `bson:"s5"`

	OffSyncType     string `bson:"offer_sync_type"`
	DirectUrl       string `bson:"direct_url"`
	SyncAdvLink     string `bson:"adv_link"`
	SyncAdvSycnCode int    `bson:"adv_sync_code"`
	SyncAdvSycnText string `bson:"adv_sync_text"`
}

type ChannelCbInfo struct {
	CbUrl  string `bson:"cburl"`
	CbCode int    `bson:"cbcode"`
	CbCnt  string `bson:"cbcnt"`
	CbErr  string `bson:"cberr"`
}

type EventLogOff struct {
	MediumId     int `bson:"medium_id"`
	OperatorId   int `bson:"operator_id"`
	AfOperatorId int `bson:"afoperator_id"`
	SalesId      int `bson:"sales_id"`
	AfSalesId    int `bson:"afsales_id"`

	NetunionId       int    `bson:"netunion_id"`
	AfNetunionId     int    `bson:"afnetunion_id"`
	OffferUniqueId   string `bson:"offer_uniqueid"`
	AfOffferUniqueId string `bson:"afoffer_uniqueid"`

	OfferId   string `bson:"offer_id"`
	AfOfferId string `bson:"afoffer_id"`

	AdvId        int `bson:"adv_id"`
	AfAdvId      int `bson:"afadv_id"`
	CreativeId   int `bson:"creative_id"`
	AfCreativeId int `bson:"afcreative_id"`
}

type EventLog struct {
	ClickId string `bson:"click_id"`

	OfferType string `bson:"offer_type"`
	OfferId   string `bson:"offer_id"`

	Eidfa        string `bson:"e_idfa"`
	Eip          string `bson:"e_ip"`
	EgoogleAdvId string `bson:"e_googleadvid"`
	EandroidId   string `bson:"e_androidid"`
	Eimei        string `bson:"e_imei"`
	EventName    string `bson:"e_name"`
	EventValue   string `bson:"e_value"`

	PayOut string `bson:"payout"`

	EventType string `bson:"event_type"`
	EventTs   int64  `bson:"event_ts"`
	EventDay  string `bson:"event_day"`
	EventHour int    `bson:"event_hour"`
	EventMin  int    `bson:"event_min"`
	EventSec  int    `bson:"event_sec"`

	ReqUrl string `bson:"request_url"`
	ReqIp  string `bson:"request_ip"`
	ReqUa  string `bson:"request_ua"`

	Offer  EventLogOff   `bson:"offer"`
	CbInfo ChannelCbInfo `bson:"callback_info"`

	ClickInfo ClickLog `bson:"click_info"`
}
