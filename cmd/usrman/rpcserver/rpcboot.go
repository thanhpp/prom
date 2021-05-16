package rpcserver

import (
	"context"
	"fmt"

	"github.com/thanhpp/prom/pkg/booting"
	"github.com/thanhpp/prom/pkg/configs"
	"github.com/thanhpp/prom/pkg/etcdclient"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
	"google.golang.org/grpc"
)

func StartGRPC(c *configs.GRPCConfig) (daemon booting.Daemon, err error) {
	logger.Get().Info("Add endpoint to etcd")
	if err := etcdclient.Get().AddEndpoints(context.Background(), c.Name, fmt.Sprintf("%s:%s", c.PublicHost, c.Port)); err != nil {
		return nil, err
	}

	daemon, err = booting.NewGRPCDaemon(c,
		func(s *grpc.Server) {
			usrmanrpc.RegisterUsrManSrvServer(s, new(usrManSrv))
		})

	if err != nil {
		return nil, err
	}

	return daemon, nil
}
