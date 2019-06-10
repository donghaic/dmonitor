package main

import (
	"destroyer-monitor/config"
	"destroyer-monitor/lib/zap"
	"destroyer-monitor/utils"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pflag.String("conf", "./config.yaml", "set configuration `file`")
	pflag.String("profile", "dev", "app profile")
	pflag.String("log-dir", "./logs", "server logs dir")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	configFile := viper.GetString("conf")
	var logger = zap.Get()
	if utils.IsEmpty(configFile) {
		pflag.Usage()
		logger.Info("can't not found config file")
		os.Exit(1)
	}

	cnf, err := config.ReadConfig(configFile)
	if err != nil {
		logger.Error("read cnf file error. file=", configFile)
		log.Fatalf("read cnf file error: %s", err)
	}

	go listenToSystemSignals()

	clickServer := NewClickServer(cnf)
	err = clickServer.Init()
	if err != nil {
		panic("click server init error")
		os.Exit(10)
	}
	clickServer.Run()

}

func listenToSystemSignals() {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
		case sig := <-signalChan:
			println(fmt.Sprintf("System signal: %s", sig))
			os.Exit(10)
		}
	}
}
