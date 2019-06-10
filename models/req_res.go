package models

// 点击事件请求信息
type ClickParams struct {
	Api          string
	OfferId      string
	ChannelId    string //渠道id
	SubChannel   string //子渠道ID
	ChannelClkId string //渠道点击ID

	GoogleAdvId string //Google Adervertiser Id
	Idfa        string //Idfa
	AndroidId   string //Android Id
	Imei        string //Android Id
	Ip          string
	Ua          string
	S1          string
	S2          string
	S3          string
	S4          string
	S5          string

	Callback string

	ClickId string //生成的点击ID

	ClickTs        int64
	RequestUrl     string
	RequestUa      string
	RequestIp      string
	RequestCountry string
	RequestOs      string

	InvalidCause string
	OffSyncType  string
	ReplaceLink  string
	SyncCode     int
	SyncCnt      string
	SyncErr      string
}

// 转化事件请求信息
type EventParams struct {
	ClickId string

	Eidfa        string
	Eip          string
	EgoogleAdvId string
	EandroidId   string
	Eimei        string
	EventName    string
	EventValue   string

	PayOut string

	EventType string
	EventTs   int64
	EventDay  string
	EventHour int
	EventMin  int
	EventSec  int

	ReqUrl string
	ReqIp  string
	ReqUa  string

	OfferId string
	ClickTs string
}

type Response struct {
	Code    int
	Content string
}
