package config

import (
	"destroyer-monitor/lib/mongo"
	"destroyer-monitor/lib/redis"
	"destroyer-monitor/lib/zap"
	"github.com/spf13/viper"
)

type Config struct {
	Port int

	Logworker logworker

	Redis redisCnf

	Mongodb mongoCnf

	Queue queue

	Httpcli HttpCli

	LogDataDir string
}

func ReadConfig(configFile string) (*Config, error) {
	var logger = zap.Get()

	logger.Info("start to read config file: ", configFile)
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("read config file error.", err)
		return nil, err
	}

	var cfg = &Config{}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		logger.Error("parse config file error.", err)
		return nil, err
	}

	// TODO check cnf?

	return cfg, err
}

type redisCnf struct {
	Entity    redis.PoolOption // 推送数据 offer,channel,blacklist
	Daycap    redis.PoolOption // 投放转化数据
	Pubsub    redis.PoolOption // Redis通知订阅
	Delaytask redis.PoolOption // 延迟队列
}

type mongoCnf struct {
	Click1     mongo.DBOption
	Click2     mongo.DBOption
	Conversion mongo.DBOption
}

type HttpCli struct {
	Timeout             int
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

type logworker struct {
	Address string
	Nodes   []string
}

type queue struct {
	LocalDataDir string
	DelayDataDir string
}
