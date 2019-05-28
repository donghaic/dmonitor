package models

type ClickLog struct {
	LogType   string `json:"log_type"`
	OfferType string `json:"offer_type"` //"customer"  "api"

	ClickId string `json:"clickid"`

	ChannelId    int    `json:"channel_id"`
	ChannelName  string `json:"channel_name"`
	SubChannel   string `json:"sub_channel"`
	ChannelClkId string `json:"channel_clkid"`

	IsSwitch bool `json:"is_switch"`

	DestroyerId   string `json:"destroyer_id"` //对外的订单id  直客是campaignid api是订单独立id
	DestroyerName string `json:"destroyer_name"`

	ApiOffer    ApiOfferStr `json:"api_offer"`
	ApiOfferCap int         `json:"api_offer_cap"`
	CusOffer    CusOfferStr `json:"customer_offer"`

	IsValid      bool   `json:"is_valid"`
	InvalidCause string `json:"invalid_cause"`

	ReqestUrl      string `json:"request_url"`
	RequestUnix    int64  `json:"requestunix"`
	RequestTime    string `json:"request_time"`
	RequestIp      string `json:"request_ip"`
	RequestCountry string `json:"request_country"`
	RequestUa      string `json:"request_ua"`
	RequestOs      string `json:"request_os"`

	Ip          string `json:"ip"`
	Ua          string `json:"ua"`
	Idfa        string `json:"idfa"`
	AndroidId   string `json:"androidid"`
	GoogleAdvId string `json:"google_advid"`
	Imei        string `json:"imei"`
	Os          string `json:"os"`
	CallBack    string `json:"callback"`
	S1          string `json:"s1"`
	S2          string `json:"s2"`
	S3          string `json:"s3"`
	S4          string `json:"s4"`
	S5          string `json:"s5"`

	OffSyncType     string `json:"offer_sync_type"`
	DirectUrl       string `json:"direct_url"`
	SyncAdvLink     string `json:"adv_link"`
	SyncAdvSycnCode int    `json:"adv_sync_code"`
	SyncAdvSycnText string `json:"adv_sync_text"`
}

type CusOfferStr struct {
	MediumId     int     `json:"medium_id"`
	MediumName   string  `json:"medium_name"`
	OperatorId   int     `json:"operator_id"`
	OperatorName string  `json:"operator_name"`
	SalesId      int     `json:"sales_id"`
	SalesName    string  `json:"sales_name"`
	AdvId        int     `json:"adv_id"`
	AdvName      string  `json:"adv_name"`
	OfferId      int     `json:"offer_id"`
	OfferName    string  `json:"offer_name"`
	CreativeId   int     `json:"creative_id"`
	CreativeName string  `json:"creative_name"`
	PackageName  string  `json:"package_name"`
	Platform     string  `json:"platform"`
	InPrice      float64 `json:"in_price"`
	InCurrency   string  `json:"in_currency"`
	OutPrice     float64 `json:"out_price"`
	OutCurrency  string  `json:"out_currency"`

	AfOperatorId   int     `json:"afoperator_id"`
	AfOperatorName string  `json:"afoperator_name"`
	AfSalesId      int     `json:"afsales_id"`
	AfSalesName    string  `json:"afsales_name"`
	AfAdvId        int     `json:"afadv_id"`
	AfAdvName      string  `json:"afadv_name"`
	AfOfferId      int     `json:"afoffer_id"`
	AfofferName    string  `json:"afoffer_name"`
	AfCreativeId   int     `json:"afcreative_id"`
	AfCreativeName string  `json:"afcreative_name"`
	AfPackageName  string  `json:"afpackage_name"`
	AfPlatform     string  `json:"afplatform"`
	AfInPrice      float64 `json:"afin_price"`
	AfInCurrency   string  `json:"afout_currency"`
}

type ApiOfferStr struct {
	MediumId       int     `json:"medium_id"`
	MediumName     string  `json:"medium_name"`
	OperatorId     int     `json:"operator_id"`
	OperatorName   string  `json:"operator_name"`
	SalesId        int     `json:"sales_id"`
	SalesName      string  `json:"sales_name"`
	NetunionId     int     `json:"netunion_id"`
	NetunionName   string  `json:"netunion_name"`
	OffferUniqueId string  `json:"offer_uniqueid"`
	OfferId        string  `json:"offer_id"`
	OfferName      string  `json:"offer_name"`
	PackageName    string  `json:"package_name"`
	Platform       string  `json:"platform"`
	Country        string  `json:"country"`
	InPrice        float64 `json:"in_price"`
	InCurrency     string  `json:"in_currency"`
	OutPrice       float64 `json:"out_price"`
	OutCurrency    string  `json:"out_currency"`

	AfOperatorId     int     `json:"afoperator_id"`
	AfOperatorName   string  `json:"afoperator_name"`
	AfSalesId        int     `json:"afsales_id"`
	AfSalesName      string  `json:"afsales_name"`
	AfNetunionId     int     `json:"afnetunion_id"`
	AfNetunionName   string  `json:"afnetunion_name"`
	AfOffferUniqueId string  `json:"afoffer_uniqueid"`
	AfOfferId        string  `json:"afoffer_id"`
	AfOfferName      string  `json:"afoffer_name"`
	AfPackageName    string  `json:"afpackage_name"`
	AfPlatform       string  `json:"afplatform"`
	AfCountry        string  `json:"afcountry"`
	AfInPrice        float64 `json:"afin_price"`
	AfInCurrency     string  `json:"afin_currency"`
	AfOutPrice       float64 `json:"afout_price"`
	AfOutCurrency    string  `json:"afout_currency"`
}
