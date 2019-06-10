package dao

import (
	"destroyer-monitor/lib/mongo"
	"destroyer-monitor/models"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

const (
	CLICK_LOG_DB   = "ClickLog"
	CONV_LOG_DB    = "EventLog"
	CONV_LOG_TABLE = "event"
)

type ConvMongoDao struct {
	mongo *mongo.MongoPool
}

func (m *ConvMongoDao) SaveConversion(conv *models.EventLog) error {

	return nil
}

func (m *ConvMongoDao) FindConvEventByClickIdAndOfferId(clickId string, offerid string) ([]models.EventLog, error) {
	query := bson.M{"click_id": clickId, "offer_id": offerid}
	println(query)
	return nil, nil
}

type ClickMongoDao struct {
	mongo *mongo.MongoPool
}

func (m *ClickMongoDao) SaveClick(click *models.ClickLogEntity) error {

	return nil
}

func (m *ClickMongoDao) FindClickEventById(clickId string, clickTsStr string) (*models.ClickLogEntity, error) {
	clickTs, _ := strconv.ParseInt(clickTsStr, 10, 64)
	collName := time.Unix(clickTs/1000000000, 0).Format("2006-01-02")
	query := bson.M{"clickid": clickId}
	println(collName, query)

	result := models.ClickLogEntity{}

	return &result, nil
}

func NewClickMongoDao(mongo *mongo.MongoPool) *ClickMongoDao {
	return &ClickMongoDao{mongo: mongo}
}

func NewConvMongoDao(mongo *mongo.MongoPool) *ConvMongoDao {
	return &ConvMongoDao{mongo: mongo}
}
