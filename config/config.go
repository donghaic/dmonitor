package config

import (
	"destroyer-monitor/lib/mongo"
	"destroyer-monitor/lib/redis"
	"destroyer-monitor/lib/zap"
	"github.com/spf13/viper"
)

type Config struct {
	Port int

	Redis redisCnf

	Mongodb mongoCnf

	Httpcli HttpCli

	TaskQueueDataDir string
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
	Entity redis.PoolOption
	Pubsub redis.PoolOption
}

type mongoCnf struct {
	Log    mongo.DBOption
	Report mongo.DBOption
}

type HttpCli struct {
	Timeout             int
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}
