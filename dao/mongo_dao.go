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

	GT_REPORT_DB    = "GtStats"
	GT_REPORT_TABLE = "report"
)

type MongoDao struct {
	clickMgoClient []*mongo.MongoPool
	convMgoClient  *mongo.MongoPool
}

func NewMongoDao(clickMgoClient []*mongo.MongoPool, convMgoClient *mongo.MongoPool) *MongoDao {
	return &MongoDao{clickMgoClient, convMgoClient}
}

func (m *MongoDao) SaveClick(conv *models.ClickLogEntity) error {

	return nil
}

func (m *MongoDao) SaveConversion(conv *models.EventLog) error {

	return nil
}

func (m *MongoDao) FindConvEventByClickIdAndOfferId(clickId string, offerid string) ([]models.EventLog, error) {
	query := bson.M{"click_id": clickId, "offer_id": offerid}
	println(query)
	return nil, nil
}

func (m *MongoDao) FindClickEventById(clickId string, clickTsStr string) (*models.ClickLogEntity, error) {
	clickTs, _ := strconv.ParseInt(clickTsStr, 10, 64)
	collName := time.Unix(clickTs/1000000000, 0).Format("2006-01-02")
	query := bson.M{"clickid": clickId}
	println(collName, query)



	result := models.ClickLogEntity{}

	return &result, nil
}
