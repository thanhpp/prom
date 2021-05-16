package boot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhpp/prom/cmd/portal/repository"

	"github.com/thanhpp/prom/cmd/portal/webserver"

	"github.com/thanhpp/prom/cmd/portal/core"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/pkg/etcdclient"
	"github.com/thanhpp/prom/pkg/logger"
)

var mainContext = context.Background()

func Boot() {
	// read config
	if err := core.SetMainConfig("dev.yml"); err != nil {
		panic(err)
	}

	// set logger
	logConfig := core.GetMainConfig().Log
	if err := logger.Set("zap", "portal", core.GetMainConfig().Environment, logConfig.Level, logConfig.Color); err != nil {
		panic(err)
	}

	// connect etcd
	logger.Get().Info("CONNECTING TO ETCD ...")
	etcdConf := core.GetMainConfig().ETCD
	if err := etcdclient.Set(&etcdConf); err != nil {
		logger.Get().Errorf("Connect etcd error: %v", err)
		panic(err)
	}

	// connect to redis
	logger.Get().Info("CONNECTING TO REDIS...")
	redisConfig := core.GetMainConfig().Redis
	if err := repository.GetRedis().Set(redisConfig); err != nil {
		logger.Get().Errorf("Connect redis error: %v", err)
		panic(err)
	}

	// init usrman service
	logger.Get().Info("CONNECTING TO USER MANAGER...")
	usrmanCtx, usrmanCtxCancel := context.WithCancel(mainContext)
	defer usrmanCtxCancel()
	if err := service.SetUsrManService(usrmanCtx, "usermanager"); err != nil {
		logger.Get().Errorf("Connect usermanager error: %v", err)
		panic(err)
	}

	// FIXME: init ccmanager

	// jwt service
	service.SetJWTSrv("SECREAT-KEY")

	// start webserver
	logger.Get().Info("STARTING HTTP SERVER...")
	webConfig := core.GetMainConfig().WebServer
	webSrvCtx, webSrvCtxCancel := context.WithCancel(mainContext)
	defer webSrvCtxCancel()
	webDaemon, err := webserver.StartHTTPServer(webSrvCtx, webConfig.Host, webConfig.Port)
	if err != nil {
		logger.Get().Errorf("Start HTTP server error: %v", err)
		panic(err)
	}

	webStart, webStop := webDaemon(webSrvCtx)
	go func() {
		if err := webStart(); err != nil {
			logger.Get().Errorf("Start HTTP server error: %v", err)
			panic(err)
		}
	}()
	defer webStop()

	// stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt) //Catch the signal
	<-quit                                                                              // Wait for the signal to come
}
