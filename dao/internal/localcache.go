package internal

import (
	"destroyer-monitor/models"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var (
	defaultExpiration = 5 * time.Minute
	cleanupInterval   = 5 * time.Minute
)

type LocalCacheService struct {
	channelCache   *gocache.Cache
	blackListCache *gocache.Cache

	apiOfferCache *gocache.Cache
	cusOfferCache *gocache.Cache

	offerDayCap    *gocache.Cache
	campaignDayCap *gocache.Cache
}

func NewLocalCache() *LocalCacheService {
	return &LocalCacheService{
		channelCache:   gocache.New(defaultExpiration, cleanupInterval),
		blackListCache: gocache.New(defaultExpiration, cleanupInterval),
		apiOfferCache:  gocache.New(defaultExpiration, cleanupInterval),
		cusOfferCache:  gocache.New(defaultExpiration, cleanupInterval),
		offerDayCap:    gocache.New(defaultExpiration, cleanupInterval),
		campaignDayCap: gocache.New(defaultExpiration, cleanupInterval),
	}
}

func (l *LocalCacheService) Update(tp, id string) {

}

func (l *LocalCacheService) GetDayCap(isapi bool, offerid string, campaignid string) (bool, map[string]int) {
	return false, nil
}

func (l *LocalCacheService) GetApiOffer(offerid string) *models.ApiOffer {
	return nil
}

func (l *LocalCacheService) GetCusOffer(offerid string) *models.CustomerOffer {
	return nil
}

func (l *LocalCacheService) GetChannel(channelid string) (*models.Channel ) {
	channel, _ := l.channelCache.Get(channelid)
	info := channel.(*models.Channel)
	return info
}

func (l *LocalCacheService) GetBlacklist(offerid string) map[string]string {
	return nil
}

func (l *LocalCacheService) SetApiOffer(offer *models.ApiOffer) {

}

func (l *LocalCacheService) SetCusOffer(offer models.CustomerOffer) {

}

func (l *LocalCacheService) SetChannel(channel models.Channel) {

}

func (l *LocalCacheService) SetBlackList(blinfo map[string]string) {

}

func (l *LocalCacheService) SetOfferCap(offerid string, cvsnow int) {

}

func (l *LocalCacheService) SetCampaignCap(isapi bool, campaignid string, cvsnow int) {

}
