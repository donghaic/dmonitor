package models

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

type Response struct {
	Code    int
	Content string
}
