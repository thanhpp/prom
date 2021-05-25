package boot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhpp/prom/cmd/noti/core"
	"github.com/thanhpp/prom/cmd/noti/repository"
	"github.com/thanhpp/prom/cmd/noti/service"
	"github.com/thanhpp/prom/cmd/noti/webserver"
	"github.com/thanhpp/prom/pkg/logger"
)

func Boot() {
	var mainCtx = context.Background()

	if err := core.SetConfig("dev.yml"); err != nil {
		panic(err)
	}

	logConf := core.GetConfig().Log
	if err := logger.Set("ZAP", core.GetConfig().ServiceName, core.GetConfig().Environment, logConf.Level, logConf.Color); err != nil {
		panic(err)
	}

	dbConf := core.GetConfig().DB
	if err := repository.Get().InitDBConnection(dbConf.GenDBDSN(), dbConf.Log); err != nil {
		panic(err)
	}

	rmq := new(service.RabbitMQService)
	if err := rmq.Connect(core.GetConfig().RabbitMQURL); err != nil {
		panic(err)
	}
	rmqCtx, rmqCtxCancel := context.WithCancel(mainCtx)
	defer rmqCtxCancel()
	rmqDaemon, err := rmq.CreateMsgDaemon(rmqCtx)
	if err != nil {
		panic(err)
	}
	rmqStart, rmqStop := rmqDaemon(rmqCtx)

	webConf := core.GetConfig().Web
	webCtx, webCtxCancel := context.WithCancel(mainCtx)
	defer webCtxCancel()
	webDaemon, err := webserver.StartHTTP(webConf.Host, webConf.Port)
	if err != nil {
		panic(err)
	}
	webStart, webStop := webDaemon(webCtx)

	go func() {
		if err := rmqStart(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := webStart(); err != nil {
			panic(err)
		}
	}()

	// stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt) //Catch the signal
	<-quit
	rmqStop()
	webStop()
}
