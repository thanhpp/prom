package booting

import (
	"context"
	"fmt"
	"net"

	"github.com/thanhpp/prom/pkg/configs"
	"google.golang.org/grpc"
)

// NewGRPCDaemon returns a daemon to start & stop grpc server
func NewGRPCDaemon(c *configs.GRPCConfig, register func(*grpc.Server)) (daemon Daemon, err error) {
	// start TCP
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", c.PublicHost, c.Port))
	if err != nil {
		return nil, err
	}

	srv := grpc.NewServer(grpc.MaxConcurrentStreams(uint32(c.MaxStream)))
	register(srv)

	daemon = func(childCtx context.Context) (start func() error, stop func()) {

		start = func() (err error) {
			if err := srv.Serve(listen); err != nil {
				return err
			}
			return nil
		}

		stop = func() {
			<-childCtx.Done()
			srv.GracefulStop()
		}

		return start, stop
	}

	return daemon, nil
}
