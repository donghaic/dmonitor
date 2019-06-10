package models

const (
	SystemException          = "system internal error"
	CacheNotHit              = "cache not hit"
	RedisQueryErr            = "redis query error"
	RedisQueryNil            = "redis query nil error"
	RedisSetErr              = "redis set error"
	RedisScanErr             = "redis search error"
	RedisDeleteErr           = "redis delete error"
	RedisValueTypeErr        = "redis value type error"
	RedisKeyHugeErr          = "redis keys to large"
	RedisConnExhaust         = "redis connect pool empty"
	LogNodeConnectError      = "log node connect error"
	LogNodeSendResponseError = "log node send response error"
	MgoQueryErr              = "mongodb query error"
	MgoSetErr                = "mongodb set error"
	MgoInsertErr             = "mongodb insert error"
	JsonParseErr             = "json parse error"
	LocalCacheFull           = "local cache full"
	UncontrollableError      = "uncontrollable error"
)

type SysLog struct {
	Date       string //2017-09-19 09:09:09
	SystemName string //系统名称
	FileName   string //文件名
	FileLine   int    //行数
	Function   string //函数
	ErrType    string //错误类型
	ErrInfo    string //错误详情

	ErrLog interface{} //发生错误时候的业务日志
}
