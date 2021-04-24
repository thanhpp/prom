package rpcserver

import (
	"github.com/thanhpp/prom/pkg/booting"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/configs"
	"google.golang.org/grpc"
)

func StartGRPC(c *configs.GRPCConfig) (daemon booting.Daemon, err error) {
	daemon, err = booting.NewGRPCDaemon(c,
		func(s *grpc.Server) {
			ccmanrpc.RegisterCCManagerServer(s, new(ccManSv))
		})

	if err != nil {
		return nil, err
	}

	return daemon, nil
}
