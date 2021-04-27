package rpcserver

import (
	"github.com/thanhpp/prom/pkg/booting"
	"github.com/thanhpp/prom/pkg/configs"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
	"google.golang.org/grpc"
)

func StartGRPC(c *configs.GRPCConfig) (daemon booting.Daemon, err error) {
	daemon, err = booting.NewGRPCDaemon(c,
		func(s *grpc.Server) {
			usrmanrpc.RegisterUsrManSrvServer(s, new(usrManSrv))
		})

	if err != nil {
		return nil, err
	}

	return daemon, nil
}
