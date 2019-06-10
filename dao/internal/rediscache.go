package internal

import (
	"destroyer-monitor/config"
	"destroyer-monitor/lib/redis"
	"destroyer-monitor/models"
	"destroyer-monitor/utils"
	"encoding/json"
	redigo "github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

type RedisCache struct {
	offerPool     *redis.ConnPool
	channelPool   *redis.ConnPool
	blacklistPool *redis.ConnPool
	daycapPool    *redis.ConnPool
}

func NewRedisCache(cnf *config.Config) (*RedisCache, error) {
	return nil, nil
}

func (r *RedisCache) GetApiOffer(netunionid, campaignid string) (models.ApiOffer, *models.SysLog) {
	key := strings.Join([]string{"netunion", "offer", "netunion", netunionid, "offer", campaignid}, ":")

	conn, err := r.offerPool.GetConnection()
	if nil != err {
		return models.ApiOffer{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   27,
			Function:   "func (r *RedisCache) GetApiOffer(netunionid, campaignid string) (models.ApiOffer, *models.SysLog)",
			ErrType:    models.RedisConnExhaust,
			ErrInfo:    err.Error(),
		}
	}
	defer conn.Close()

	mapStr, err := redigo.StringMap(conn.Do("HGETALL", key))
	if nil == err {
		if _, ok := mapStr["unique_id"]; !ok {
			return models.ApiOffer{}, &models.SysLog{
				SystemName: "DestroyerClick",
				FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
				FileLine:   41,
				Function:   "func (r *RedisCache) GetRedisCampaign(id int) (models.Campaign, *models.SysLog)",
				ErrType:    models.RedisQueryNil,
				ErrInfo:    key,
			}
		}

		data := models.ApiOffer{
			OperatorName: mapStr["operator_name"],
			SalesName:    mapStr["sales_name"],
			UniqueId:     mapStr["unique_id"],
			OfferId:      mapStr["offer_id"],
			Name:         mapStr["name"],
			PackageName:  mapStr["package_name"],
			PreviewUrl:   mapStr["preview_url"],
			// "IsIncent":     mapStr["incent"],
			MinOsv: mapStr["min_osv"],
			// "InPrice":      mapStr["price"],
			InCurrency: "USD", // mapStr["in_currency"],
			// "DayCap":       mapStr["day_cap"],
			// "Countries":    mapStr["countries"],
			Platform:     mapStr["platform"],
			NetTrackLink: mapStr["net_track_link"],
			TrackLink:    mapStr["track_link"],
			// NetunionId:   mapStr["netunion_id"],
			NetunionName: mapStr["netunion_name"],
			// "Updated":      mapStr["push_ns"],
		}

		salesid, err := strconv.ParseInt(mapStr["sales_id"], 10, 32)
		if nil == err {
			data.SalesId = int(salesid)
		}
		operid, err := strconv.ParseInt(mapStr["operator_id"], 10, 32)
		if nil == err {
			data.OperatorId = int(operid)
		}
		netid, err := strconv.ParseInt(mapStr["netunion_id"], 10, 32)
		if nil == err {
			data.NetunionId = int(netid)
		}

		incent, err := strconv.ParseInt(mapStr["incent"], 10, 32)
		if nil == err {
			data.IsIncent = int(incent)
		}

		price, err := strconv.ParseFloat(mapStr["price"], 32)
		if nil == err {
			data.InPrice = price
		}

		day_cap, err := strconv.ParseInt(mapStr["day_cap"], 10, 32)
		if nil == err {
			data.DayCap = int(day_cap)
		}

		data.Countries = strings.Split(mapStr["countries"], ",")

		updated, err := strconv.ParseInt(mapStr["push_ns"], 10, 64)
		if nil == err {
			data.Updated = updated
		}

		return data, nil
	}

	if redigo.ErrNil == err {
		return models.ApiOffer{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   41,
			Function:   "func (r *RedisCache) GetRedisCampaign(id int) (models.Campaign, *models.SysLog)",
			ErrType:    models.RedisQueryNil,
			ErrInfo:    key,
		}
	}

	return models.ApiOffer{}, &models.SysLog{
		SystemName: "DestroyerClick",
		FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
		FileLine:   70,
		Function:   "func (r *RedisCache) GetRedisCampaign(id int) (models.Campaign, *models.SysLog)",
		ErrType:    models.RedisQueryErr,
		ErrInfo:    key,
	}
}

func (r *RedisCache) GetCusOffer(campaignid string) (models.CustomerOffer, *models.SysLog) {
	key := strings.Join([]string{"customer", "offer", campaignid}, ":")

	conn, err := r.offerPool.GetConnection()
	if nil != err {
		return models.CustomerOffer{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   27,
			Function:   "func (r *RedisCache) GetApiOffer(netunionid, campaignid string) (models.ApiOffer, *models.SysLog)",
			ErrType:    models.RedisConnExhaust,
			ErrInfo:    err.Error(),
		}
	}
	defer conn.Close()

	jsonStr, err := redigo.String(conn.Do("GET", key))
	if nil == err {
		data := models.CustomerOffer{}
		err = json.Unmarshal(utils.StrToBytes(jsonStr), &data)
		if nil == err {
			return data, nil
		}
		return models.CustomerOffer{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   41,
			Function:   "func (r *RedisCache) GetRedisCampaign(id int) (models.Campaign, *models.SysLog)",
			ErrType:    models.JsonParseErr,
			ErrInfo:    jsonStr,
		}
	}

	if redigo.ErrNil == err {
		return models.CustomerOffer{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   41,
			Function:   "func (r *RedisCache) GetRedisCampaign(id int) (models.Campaign, *models.SysLog)",
			ErrType:    models.RedisQueryNil,
			ErrInfo:    key,
		}
	}

	return models.CustomerOffer{}, &models.SysLog{
		SystemName: "DestroyerClick",
		FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
		FileLine:   70,
		Function:   "func (r *RedisCache) GetRedisCampaign(id int) (models.Campaign, *models.SysLog)",
		ErrType:    models.RedisQueryErr,
		ErrInfo:    key,
	}
}

func (r *RedisCache) GetRedisChannel(channelid string) (models.Channel, *models.SysLog) {
	key := strings.Join([]string{"channel", "info", "channel", channelid}, ":")

	conn, err := r.channelPool.GetConnection()
	if nil != err {
		return models.Channel{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   27,
			Function:   "func (r *RedisCache) GetRedisChannel(id int) (models.Campaign, *models.SysLog)",
			ErrType:    models.RedisConnExhaust,
			ErrInfo:    err.Error(),
		}
	}
	defer conn.Close()

	jsonStr, err := redigo.StringMap(conn.Do("HGETALL", key))
	if nil == err {
		if _, ok := jsonStr["id"]; !ok {
			return models.Channel{}, &models.SysLog{
				SystemName: "DestroyerClick",
				FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
				FileLine:   41,
				Function:   "func (r *RedisCache) GetRedisChannel(id int) (models.Campaign, *models.SysLog)",
				ErrType:    models.JsonParseErr,
				ErrInfo:    key,
			}
		}

		data := models.Channel{
			// "Id":            jsonStr["id"],
			Name:      jsonStr["name"],
			Callback:  jsonStr["callback_url"],
			MediaName: jsonStr["media_name"],
			// "PriceDiscount": jsonStr["price_discount"],
			// "DeductionPer":  jsonStr["deduction_percentage"],
		}
		id, err := strconv.ParseInt(jsonStr["id"], 10, 32)
		if nil == err {
			data.Id = int(id)
		}

		mediaid, err := strconv.ParseInt(jsonStr["media"], 10, 32)
		if nil == err {
			data.Media = int(mediaid)
		}

		priced, err := strconv.ParseInt(jsonStr["price_discount"], 10, 32)
		if nil == err {
			data.PriceDiscount = int(priced)
		}

		cbdis, err := strconv.ParseInt(jsonStr["deduction_percentage"], 10, 32)
		if nil == err {
			data.DeductionPer = int(cbdis)
		}

		return data, nil
	}

	if redigo.ErrNil == err {
		return models.Channel{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   41,
			Function:   "func (r *RedisCache) GetRedisChannel(id int) (models.Campaign, *models.SysLog)",
			ErrType:    models.RedisQueryNil,
			ErrInfo:    key,
		}
	}

	return models.Channel{}, &models.SysLog{
		SystemName: "DestroyerClick",
		FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
		FileLine:   41,
		Function:   "func (r *RedisCache) GetRedisChannel(id int) (models.Campaign, *models.SysLog)",
		ErrType:    models.RedisQueryErr,
		ErrInfo:    key,
	}
}

func (r *RedisCache) GetRedisBlackList(campaignid string) (map[string]string, *models.SysLog) {
	key := strings.Join([]string{"blacklist", "offer", campaignid}, ":")

	conn, err := r.blacklistPool.GetConnection()
	if nil != err {
		return map[string]string{}, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   27,
			Function:   "func (r *RedisCache) GetRedisBlackList(id int) (models.Campaign, *models.SysLog)",
			ErrType:    models.RedisConnExhaust,
			ErrInfo:    err.Error(),
		}
	}
	defer conn.Close()

	jsonStr, err := redigo.StringMap(conn.Do("HGETALL", key))
	if nil == err {
		return jsonStr, nil
	}

	if redigo.ErrNil == err {
		return map[string]string{}, nil
	}

	return map[string]string{}, &models.SysLog{
		SystemName: "DestroyerClick",
		FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
		FileLine:   41,
		Function:   "func (r *RedisCache) GetRedisBlackList(id int) (models.Campaign, *models.SysLog)",
		ErrType:    models.RedisQueryErr,
		ErrInfo:    key,
	}
}

func (r *RedisCache) GetOfferDayStats(offerid string) (int, *models.SysLog) {
	filedt := time.Now().Format("2006-01-02")

	conn, err := r.daycapPool.GetConnection()
	if nil != err {
		return 0, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   87,
			Function:   "func (r *RedisCache) GetOfferCampaignStats(offerid string) (models.DayStatsCap, *models.SysLog)",
			ErrType:    models.RedisQueryErr,
			ErrInfo:    err.Error(),
		}
	}
	defer conn.Close()

	camcap, err := redigo.Int(conn.Do("HGET", strings.Join([]string{"offer", offerid}, ":"), filedt))
	if nil != err && redigo.ErrNil != err {
		return 0, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   106,
			Function:   "func (r *RedisCache) GetOfferCampaignStats(offerid, campaignid ing) (models.DayStatsCap, *models.SysLog)",
			ErrType:    models.RedisQueryErr,
			ErrInfo:    err.Error(),
		}
	}

	return camcap, nil
}

func (r *RedisCache) GetCampaignDayStats(isapi bool, campaignid string) (int, *models.SysLog) {
	key := strings.Join([]string{"cus", "campaignid", campaignid}, ":")
	if isapi {
		key = strings.Join([]string{"api", "campaignid", campaignid}, ":")
	}

	filedt := time.Now().Format("2006-01-02")

	conn, err := r.daycapPool.GetConnection()
	if nil != err {
		return 0, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   87,
			Function:   "func (r *RedisCache) GetOfferCampaignStats(offerid string) (models.DayStatsCap, *models.SysLog)",
			ErrType:    models.RedisQueryErr,
			ErrInfo:    err.Error(),
		}
	}
	defer conn.Close()

	camcap, err := redigo.Int(conn.Do("HGET", key, filedt))
	if nil != err && redigo.ErrNil != err {
		return 0, &models.SysLog{
			SystemName: "DestroyerClick",
			FileName:   "src/mcad.com/OnlineClk/MemCache/RedisCache.go",
			FileLine:   106,
			Function:   "func (r *RedisCache) GetOfferCampaignStats(offerid, campaignid ing) (models.DayStatsCap, *models.SysLog)",
			ErrType:    models.RedisQueryErr,
			ErrInfo:    err.Error(),
		}
	}

	return camcap, nil
}

func (r *RedisCache) SetSmartLink(idsstr string) error {
	conn, err := r.offerPool.GetConnection()
	if nil != err {
		log.Println("SetSmartLink. get redis connection error.", err.Error())
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", "smartlink:lige", idsstr)
	if err != nil {
		log.Println("SetSmartLink. set smartlink:lige error.", err.Error())
		return err
	}

	return nil
}

func (r *RedisCache) GetSmartLink() string {
	conn, err := r.offerPool.GetConnection()
	if nil != err {
		log.Println("SetSmartLink. get redis connection error.", err.Error())
		return ""
	}
	defer conn.Close()

	dstr, err := redigo.String(conn.Do("GET", "smartlink:lige"))
	if err != nil {
		log.Println("SetSmartLink. get smartlink:lige error.", err.Error())
		return ""
	}
	return dstr
}
