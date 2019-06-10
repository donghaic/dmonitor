package api

import (
	"destroyer-monitor/dao"
	"destroyer-monitor/services"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type SmartLinkApi struct {
	redisDao *dao.RedisDao
}

func (s *SmartLinkApi) GetSmartLink(ctx *fasthttp.RequestCtx) {
	link, _ := s.redisDao.GetSmartLink()
	resdata := map[string]interface{}{
		"smartlink_offerids": services.SmarklinkOffers,
		"slredis_offerids":   link,
	}
	resbyte, _ := json.Marshal(resdata)
	_, _ = ctx.WriteString(string(resbyte))
}

func (s *SmartLinkApi) SetSmartLink(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	if 0 < len(body) {
		err := s.redisDao.SetSmartLink(string(body))
		if err != nil {
			_, _ = ctx.WriteString(err.Error())
			return
		}
		_, _ = ctx.WriteString("okay")
		return
	}
	_, _ = ctx.WriteString("empty!")
}
