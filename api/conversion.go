package api

import (
	"destroyer-monitor/models"
	"destroyer-monitor/services/handler"
	"destroyer-monitor/utils"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Conversion struct {
	handler *handler.ConvHandler
}

func NewConversion(handler *handler.ConvHandler) *Conversion {
	return &Conversion{handler,}
}

func (c *Conversion) Handle(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()
	params := &models.EventParams{
		ClickId:      string(args.Peek("click_id")),
		Eidfa:        string(args.Peek("e_idfa")),
		Eip:          string(args.Peek("e_ip")),
		EgoogleAdvId: string(args.Peek("e_googleadvid")),
		EandroidId:   string(args.Peek("e_androidid")),
		Eimei:        string(args.Peek("e_imei")),
		EventName:    string(args.Peek("e_name")),
		EventValue:   string(args.Peek("e_value")),
		PayOut:       string(args.Peek("payout")),
	}

	if 0 == len(params.ClickId) {
		params.ClickId = string(args.Peek("clickid"))
	}

	if 0 == len(params.EventName) {
		params.EventName = "fopen"
	}

	if 0 == len(params.ClickId) {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		_, _ = ctx.WriteString("NotFound")
		return
	}

	clkidsplit := strings.Split(params.ClickId, "r")
	clkidsplitLen := len(clkidsplit)

	if clkidsplitLen < 3 {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		_, _ = ctx.WriteString("NotFound")
		return
	}

	if 3 < clkidsplitLen {
		_, err := strconv.ParseInt(clkidsplit[clkidsplitLen-2], 10, 64)
		if nil != err {
			ctx.Response.SetStatusCode(http.StatusNotFound)
			_, _ = ctx.WriteString("NotFound")
			return
		}
	}

	offerid := strings.Join(clkidsplit[:clkidsplitLen-2], "r")
	if 0 != strings.Compare("wy", string([]rune(offerid)[:2])) {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		_, _ = ctx.WriteString("NotFound")
		return
	}

	params.OfferId = string([]rune(offerid)[2:])
	params.ClickTs = clkidsplit[clkidsplitLen-2]

	params.EventType = "normal"
	params.EventTs = time.Now().UnixNano()
	params.EventDay = time.Now().Format("2006-01-02")
	params.EventHour = time.Now().Hour()
	params.EventMin = time.Now().Minute()
	params.EventSec = time.Now().Second()
	params.ReqUrl = string(ctx.RequestURI())
	params.ReqIp = utils.FastGetIp(ctx)
	params.ReqUa = string(ctx.UserAgent())

	res := c.handler.Handle(params)

	//返回结果
	if 302 == res.Code {
		ctx.Redirect(res.Content, 302)
	} else {
		ctx.SetStatusCode(res.Code)
		_, _ = ctx.WriteString(res.Content)
	}
}
