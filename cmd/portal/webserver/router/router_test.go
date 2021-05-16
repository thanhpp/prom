package router_test

import (
	"context"
	"testing"

	"github.com/thanhpp/prom/cmd/portal/repository/redisdb"

	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/router"
)

func TestSetUsrManService(t *testing.T) {

}

func TestNewRouter(t *testing.T) {
	// usrman service
	var (
		ctx    = context.Background()
		target = "127.0.0.1:8090"
	)

	if err := service.SetUsrManService(ctx, target); err != nil {
		t.Error(err)
		return
	}

	// jwt service
	var (
		key = "key"
	)
	service.SetJWTSrv(key)

	// redis service
	var (
		conf = redisdb.RedisConfig{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}
	)

	if err := redisdb.Set(conf); err != nil {
		t.Error(err)
		return
	}

	// start router
	r := router.NewRouter()
	if err := r.Run(":12345"); err != nil {
		t.Error(err)
		return
	}
}
