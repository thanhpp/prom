package rpcserver

import (
	"context"
	"fmt"

	"github.com/thanhpp/prom/pkg/booting"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/configs"
	"github.com/thanhpp/prom/pkg/etcdclient"
	"google.golang.org/grpc"
)

func StartGRPC(c *configs.GRPCConfig, shardID int64) (daemon booting.Daemon, err error) {
	if c.DockerMode {
		if err := etcdclient.Get().AddEndpoints(context.Background(), fmt.Sprintf("%s-%d", c.Name, shardID), fmt.Sprintf("%s:%s", c.Name, c.Port)); err != nil {
			return nil, err
		}
	} else {
		if err := etcdclient.Get().AddEndpoints(context.Background(), fmt.Sprintf("%s-%d", c.Name, shardID), fmt.Sprintf("%s:%s", c.PublicHost, c.Port)); err != nil {
			return nil, err
		}
	}

	daemon, err = booting.NewGRPCDaemon(c,
		func(s *grpc.Server) {
			ccmanrpc.RegisterCCManagerServer(s, new(ccManSv))
		})

	if err != nil {
		return nil, err
	}

	return daemon, nil
}
