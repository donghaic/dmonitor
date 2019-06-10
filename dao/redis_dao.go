package dao

import (
	"destroyer-monitor/config"
	"destroyer-monitor/dao/internal"
	"destroyer-monitor/models"
)

type CacheType string

const (
	Offer_Cache CacheType = "offer"
)

type RedisDao struct {
	localCache *internal.LocalCacheService
	redisCache *internal.RedisCache
}

func NewRedisDao(conf *config.Config) (*RedisDao, error) {
	return nil, nil
}

func (m *RedisDao) GetApiOffer(netunionid, campaignid string) (*models.ApiOffer, *models.SysLog) {

	offcache := m.localCache.GetApiOffer(campaignid)
	if nil != offcache {
		return offcache, nil
	}

	offer, err := m.redisCache.GetApiOffer(netunionid, campaignid)
	if nil == err {
		go m.localCache.SetApiOffer(offer)
		return &offer, nil
	}

	return nil, err
}

func (m *RedisDao) GetCusOffer(campaignid string) (*models.CustomerOffer, *models.SysLog) {

	offcache := m.localCache.GetCusOffer(campaignid)
	if nil != offcache {
		return offcache, nil
	}

	offer, err := m.redisCache.GetCusOffer(campaignid)
	if nil == err {
		go m.localCache.SetCusOffer(offer)

		return &offer, nil
	}

	return nil, err
}

func (m *RedisDao) GetChannel(channelid string) (*models.Channel, *models.SysLog) {

	offcache := m.localCache.GetChannel(channelid)
	if nil != offcache {
		return offcache, nil
	}

	offer, err := m.redisCache.GetRedisChannel(channelid)
	if nil == err {
		go m.localCache.SetChannel(offer)

		return &offer, nil
	}

	return nil, err
}

func (m *RedisDao) GetBlacklist(campaignid string) (map[string]string, *models.SysLog) {

	offcache := m.localCache.GetBlacklist(campaignid)
	if _, ok := offcache["type"]; ok {
		return offcache, nil
	}

	offer, err := m.redisCache.GetRedisBlackList(campaignid)
	if nil == err {
		go m.localCache.SetBlackList(offer)

		return offer, nil
	}

	return map[string]string{}, err
}

func (m *RedisDao) GetDayCap(isapi bool, offerid, campaignid string) (map[string]int, *models.SysLog) {
	if isapi {

		isupdate, offcache := m.localCache.GetDayCapInfo(true, "", campaignid)
		if nil != offcache && !isupdate {
			return offcache, nil
		}

		offer, err := m.redisCache.GetCampaignDayStats(true, campaignid)
		if nil == err {
			go m.localCache.SetCampaignCap(true, campaignid, offer)

			return map[string]int{
				"offer":    0,
				"campaign": offer,
			}, nil
		}

		return nil, err

	} else {

		isupdate, offcache := m.localCache.GetDayCapInfo(false, offerid, campaignid)
		if nil != offcache && !isupdate {
			return offcache, nil
		}

		offercvs, _ := m.redisCache.GetOfferDayStats(offerid)
		camcvs, err := m.redisCache.GetCampaignDayStats(false, campaignid)
		if nil == err {
			go m.localCache.SetCampaignCap(false, campaignid, camcvs)
			go m.localCache.SetOfferCap(offerid, offercvs)

			return map[string]int{
				"offer":    offercvs,
				"campaign": camcvs,
			}, nil
		}

		return nil, err
	}
}

func (m *RedisDao) Update(s string, s2 string) {

}
func (r *RedisDao) GetOfferDayStats(offerid string) (int, *models.SysLog) {

	return 0, nil
}

func (r *RedisDao) GetCampaignDayStats(isapi bool, campaignid string) (int, *models.SysLog) {

	return 0, nil
}

func (r *RedisDao) SetSmartLink(idsstr string) error {

	return nil
}

func (r *RedisDao) GetSmartLink() string {

	return ""
}
