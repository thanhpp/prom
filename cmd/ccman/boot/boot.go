package boot

import (
	"context"

	"github.com/thanhpp/prom/pkg/etcdclient"

	"github.com/thanhpp/prom/cmd/ccman/repository"

	"github.com/thanhpp/prom/cmd/ccman/core"
	"github.com/thanhpp/prom/cmd/ccman/rpcserver"
	"github.com/thanhpp/prom/pkg/logger"
)

func Boot(shardID int64) (err error) {
	var (
		ctx = context.Background()
	)

	if err := core.SetMainConfig("dev.yml"); err != nil {
		return err
	}
	if shardID > 0 {
		core.GetConfig().ShardID = shardID
	}

	logConfig := core.GetConfig().Log
	if err := logger.Set("ZAP", "ccmanager", "DEVELOPMENT", logConfig.Level, logConfig.Color); err != nil {
		return err
	}

	logger.Get().Info("CONNECTING TO DB")
	if err := repository.GetDAO().InitDBConnection(core.GetConfig().DB.GenDBDSN(), core.GetConfig().DB.Log); err != nil {
		return err
	}
	if err := repository.GetDAO().AutoMigrate(ctx); err != nil {
		return err
	}

	logger.Get().Info("CONENCTING TO ETCD")
	etcdConf := &core.GetConfig().ETCD
	if err := etcdclient.Set(etcdConf); err != nil {
		return err
	}

	logger.Get().Info("STARTING GRPC SERVER")
	gRPCdaemon, err := rpcserver.StartGRPC(&core.GetConfig().GRPC, core.GetConfig().ShardID)
	if err != nil {
		return err
	}
	start, stop := gRPCdaemon(ctx)
	go func() {
		if err := start(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	stop()
	return nil
}
