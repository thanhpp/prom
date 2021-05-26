package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thanhpp/prom/cmd/noti/webserver/router"
	"github.com/thanhpp/prom/pkg/booting"
	"github.com/thanhpp/prom/pkg/logger"
)

func StartHTTP(host string, port string) (d booting.Daemon, err error) {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: router.NewRouter(),
	}

	d = func(ctx context.Context) (start func() error, stop func()) {
		start = func() error {
			logger.Get().Info("Starting HTTP Server...")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Get().Errorf("Server shutdown: %v", err)
				return err
			}

			return nil
		}

		stop = func() {
			logger.Get().Info("SHUTTING DOWN SERVER .....")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Fatalf("SHUTDOWN SERVER - FORCED: %+v\n", err)
			}
		}
		return start, stop
	}
	return d, nil
}
