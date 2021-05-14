package repository

import (
	"context"
	"time"

	"github.com/thanhpp/prom/cmd/portal/repository/redisdb"
)

type iDao interface {
	SetKey(ctx context.Context, key string, value string, timeout time.Duration) (err error)
	GetValue(ctx context.Context, key string) (value string, err error)
	DeleteKey(ctx context.Context, key string) (err error)
}

func GetRedis() iDao {
	return redisdb.Get()
}
