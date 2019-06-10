package macro

import (
	"destroyer-monitor/models"
	"destroyer-monitor/utils"
	"strconv"
	"strings"
)

func ReplacedClickTLMacroAndFunc(url string, eveParams *models.EventParams, clklog *models.ClickLogEntity) string {
	relink := url

	relink = strings.Replace(relink, "{clickid}", clklog.ChannelClkId, -1)
	relink = strings.Replace(relink, "{aff_sub}", clklog.ChannelClkId, -1)
	relink = strings.Replace(relink, "{idfa}", clklog.Idfa, -1)
	relink = strings.Replace(relink, "{gaid}", clklog.GoogleAdvId, -1)
	relink = strings.Replace(relink, "{android_id}", clklog.AndroidId, -1)
	relink = strings.Replace(relink, "{subchannel}", clklog.SubChannel, -1)
	relink = strings.Replace(relink, "{channel}", clklog.SubChannel, -1)
	relink = strings.Replace(relink, "{s1}", clklog.S1, -1)
	relink = strings.Replace(relink, "{aff_sub2}", clklog.S1, -1)
	relink = strings.Replace(relink, "{s2}", clklog.S2, -1)
	relink = strings.Replace(relink, "{s3}", clklog.S3, -1)
	relink = strings.Replace(relink, "{s4}", clklog.S4, -1)
	relink = strings.Replace(relink, "{s5}", clklog.S5, -1)
	if 0 != clklog.ApiOffer.NetunionId {
		relink = strings.Replace(relink, "{payout}", strconv.FormatFloat(utils.Decimal(float64(clklog.ApiOffer.OutPrice)), 'G', -1, 32), -1)
	} else {
		relink = strings.Replace(relink, "{payout}", strconv.FormatFloat(utils.Decimal(float64(clklog.CusOffer.OutPrice)), 'G', -1, 32), -1)
	}

	return relink
}
