package boot

import (
	"context"

	"github.com/thanhpp/prom/cmd/ccman/core"
	"github.com/thanhpp/prom/cmd/ccman/rpcserver"
	"github.com/thanhpp/prom/pkg/logger"
)

func Boot() (err error) {
	var (
		ctx = context.Background()
	)

	if err := core.SetMainConfig("dev.yml"); err != nil {
		return err
	}

	logConfig := core.GetConfig().Log
	if err := logger.Set("ZAP", "ccmanager", "DEVELOPMENT", logConfig.Level, logConfig.Color); err != nil {
		return err
	}

	logger.Get().Info("STARTING GRPC SERVER")
	gRPCdaemon, err := rpcserver.StartGRPC(&core.GetConfig().GRPC)
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
