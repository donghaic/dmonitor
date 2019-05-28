package zap

import (
	"github.com/arthurkiller/rollingwriter"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var serverLogger *zap.SugaredLogger = nil
var eventLogger *zap.Logger = nil
var serverOnce, eventOnce sync.Once

func init() {
	viper.SetDefault("log-dir", "./logs")
	viper.SetDefault("server-log", "server")
	viper.SetDefault("access-log", "access")
	viper.SetDefault("profile", "dev")
}

func Get() *zap.SugaredLogger {
	if serverLogger == nil {
		newSugarLogger()
	}
	return serverLogger
}

func GetEvent() *zap.Logger {
	if eventLogger == nil {
		newEventLogger()
	}
	return eventLogger
}

func newSugarLogger() *zap.SugaredLogger {
	serverOnce.Do(func() {
		profile := viper.GetString("profile")
		logDir := viper.GetString("log-dir")
		logName := viper.GetString("server-log")
		config := zap.NewDevelopmentEncoderConfig()
		l := doCreate(profile, logDir, logName, config)
		serverLogger = l.Sugar()
		serverLogger.Info("newSugarLogger done", )
	})

	return serverLogger

}

func newEventLogger() *zap.Logger {
	eventOnce.Do(func() {
		profile := viper.GetString("profile")
		logDir := viper.GetString("log-dir")
		logName := viper.GetString("access-log")

		var config = zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			MessageKey:     "M",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		l := doCreate(profile, logDir, logName, config)
		eventLogger = l
		serverLogger.Info("newEventLogger done")
	})

	return eventLogger

}

func doCreate(profile, logDir, logName string, cnf zapcore.EncoderConfig) *zap.Logger {

	var logger *zap.Logger
	if profile == "prod" {
		// rolling 配置
		config := rollingwriter.Config{
			LogPath:       logDir,
			FileName:      logName,
			TimeTagFormat: "20060102150405",
			MaxRemain:     1,
			//			RollingTimePattern: "0 0 0 * * *", // Rolling at 00:00 AM everyday
			RollingVolumeSize: "1G",
			WriterMode:        "async",
		}

		// 创建 writer
		w, _ := rollingwriter.NewWriterFromConfig(&config)

		// 创建 logger
		ws := zapcore.AddSync(w)

		encoder := zapcore.NewConsoleEncoder(cnf)
		core := zapcore.NewCore(
			encoder,
			ws,
			zap.InfoLevel,
		)

		logger = zap.New(core)
	} else {
		var cfg zap.Config
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		l, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		logger = l

	}

	return logger
}
