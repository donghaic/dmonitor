package api

import (
	"destroyer-monitor/models"
	"destroyer-monitor/services/handler"
	"destroyer-monitor/utils"
	"errors"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type ClickApi struct {
	handler *handler.ClickHandler
}

func NewClick(handler *handler.ClickHandler) *ClickApi {
	return &ClickApi{handler,}
}

func (c *ClickApi) Handle(ctx *fasthttp.RequestCtx) {
	clickParams, err := c.parseReq(ctx)
	if err != nil {
		ctx.SetStatusCode(404)
		_, _ = ctx.WriteString("Required parameter missing")
		return
	}
	res := c.handler.Handle(clickParams)
	//返回结果
	if 302 == res.Code {
		ctx.Redirect(res.Content, 302)
	} else {
		ctx.SetStatusCode(res.Code)
		_, _ = ctx.WriteString(res.Content)
	}
}

func (c *ClickApi) parseReq(ctx *fasthttp.RequestCtx) (*models.ClickParams, error) {
	args := ctx.URI().QueryArgs()
	clkParam := models.ClickParams{
		Api:          string(args.Peek("api")),
		OfferId:      string(args.Peek("offer_id")),
		ChannelId:    string(args.Peek("channel_id")),
		SubChannel:   string(args.Peek("subchannel")),
		ChannelClkId: string(args.Peek("channel_click_id")),
		GoogleAdvId:  string(args.Peek("google_advertiser_id")),
		Idfa:         string(args.Peek("idfa")),
		AndroidId:    string(args.Peek("androidid")),
		Imei:         string(args.Peek("imei")),
		Ip:           string(args.Peek("ip")),
		Ua:           string(args.Peek("ua")),
		S1:           string(args.Peek("s1")),
		S2:           string(args.Peek("s2")),
		S3:           string(args.Peek("s3")),
		S4:           string(args.Peek("s4")),
		S5:           string(args.Peek("s5")),
		Callback:     string(args.Peek("callback")),
	}
	if 0 < len(string(args.Peek("aff_sub"))) && !strings.Contains(string(args.Peek("aff_sub")), "aff_sub") {
		clkParam.ChannelClkId = string(args.Peek("aff_sub"))
	}
	if 0 < len(string(args.Peek("channel"))) && !strings.Contains(string(args.Peek("channel")), "channel") {
		clkParam.SubChannel = string(args.Peek("channel"))
	}
	if 0 < len(string(args.Peek("gaid"))) && !strings.Contains(string(args.Peek("gaid")), "gaid") {
		clkParam.GoogleAdvId = string(args.Peek("gaid"))
	}
	if 0 < len(string(args.Peek("aff_sub2"))) && !strings.Contains(string(args.Peek("aff_sub2")), "aff_sub2") {
		clkParam.S1 = string(args.Peek("aff_sub2"))
	}
	clkParam.RequestUa = string(ctx.UserAgent())
	clkParam.RequestIp = utils.FastGetIp(ctx)
	clkParam.RequestUrl = string(ctx.RequestURI())
	clkParam.ClickTs = time.Now().UnixNano()
	clkParam.RequestCountry = ""
	clkParam.RequestOs = ""
	if 0 == len(clkParam.Api) || 0 == len(clkParam.OfferId) || 0 == len(clkParam.ChannelId) {
		return &clkParam, errors.New("required parameter missing")
	}
	uapara := utils.DInfoFUa{}
	if 0 < len(clkParam.Ua) {
		uapara = utils.UserAgentParse(clkParam.Ua)
	} else {
		uapara = utils.UserAgentParse(clkParam.RequestUa)
	}
	clkParam.RequestOs = uapara.Os

	return &clkParam, nil
}
