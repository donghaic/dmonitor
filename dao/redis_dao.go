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

func (r *RedisDao) GetApiOffer(netunionid, offerId string) (*models.ApiOffer, *models.SysLog) {

	localOffer := r.localCache.GetApiOffer(offerId)
	if nil != localOffer {
		return localOffer, nil
	}

	offer, err := r.redisCache.GetApiOffer(netunionid, offerId)
	if nil == err {
		go r.localCache.SetApiOffer(offer)
		return offer, nil
	}

	return nil, err
}

func (r *RedisDao) GetCusOffer(campaignid string) (*models.CustomerOffer, *models.SysLog) {

	offcache := r.localCache.GetCusOffer(campaignid)
	if nil != offcache {
		return offcache, nil
	}

	offer, err := r.redisCache.GetCusOffer(campaignid)
	if nil == err {
		go r.localCache.SetCusOffer(offer)

		return &offer, nil
	}

	return nil, err
}

func (r *RedisDao) GetChannel(channelid string) (*models.Channel, *models.SysLog) {

	offcache := r.localCache.GetChannel(channelid)
	if nil != offcache {
		return offcache, nil
	}

	offer, err := r.redisCache.GetRedisChannel(channelid)
	if nil == err {
		go r.localCache.SetChannel(offer)
		return &offer, nil
	}
	return nil, err
}

func (r *RedisDao) GetBlacklist(campaignid string) (map[string]string, *models.SysLog) {

	offcache := r.localCache.GetBlacklist(campaignid)
	if _, ok := offcache["type"]; ok {
		return offcache, nil
	}

	offer, err := r.redisCache.GetRedisBlackList(campaignid)
	if nil == err {
		go r.localCache.SetBlackList(offer)

		return offer, nil
	}

	return map[string]string{}, err
}

func (r *RedisDao) GetDayCap(isapi bool, offerid, campaignid string) (map[string]int, *models.SysLog) {
	if isapi {

		isupdate, offcache := r.localCache.GetDayCap(true, "", campaignid)
		if nil != offcache && !isupdate {
			return offcache, nil
		}

		offer, err := r.redisCache.GetCampaignDayStats(true, campaignid)
		if nil == err {
			go r.localCache.SetCampaignCap(true, campaignid, offer)

			return map[string]int{
				"offer":    0,
				"campaign": offer,
			}, nil
		}

		return nil, err

	} else {

		isupdate, offcache := r.localCache.GetDayCap(false, offerid, campaignid)
		if nil != offcache && !isupdate {
			return offcache, nil
		}

		offercvs, _ := r.redisCache.GetOfferDayStats(offerid)
		camcvs, err := r.redisCache.GetCampaignDayStats(false, campaignid)
		if nil == err {
			go r.localCache.SetCampaignCap(false, campaignid, camcvs)
			go r.localCache.SetOfferCap(offerid, offercvs)

			return map[string]int{
				"offer":    offercvs,
				"campaign": camcvs,
			}, nil
		}

		return nil, err
	}
}

func (r *RedisDao) Update(s string, s2 string) {

}
func (r *RedisDao) GetOfferDayStats(offerid string) (int, *models.SysLog) {

	return 0, nil
}

func (r *RedisDao) GetCampaignDayStats(isapi bool, campaignid string) (int, *models.SysLog) {

	return 0, nil
}

func (r *RedisDao) SetSmartLink(idsstr string) error {
	err := r.redisCache.SetSmartLink(idsstr)
	return err
}

func (r *RedisDao) GetSmartLink() (string, error) {
	smartLink, err := r.redisCache.GetSmartLink()
	return smartLink, err
}
