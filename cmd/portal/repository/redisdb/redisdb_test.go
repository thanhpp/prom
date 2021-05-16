package redisdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/thanhpp/prom/cmd/portal/repository/redisdb"
)

func TestSetRDB(t *testing.T) {
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
}

func TestSetKey(t *testing.T) {
	TestSetRDB(t)
	var (
		ctx     = context.Background()
		key     = "4176ae56-e924-445c-9407-007eca2cae7f"
		value   = "ok"
		timeout = time.Minute
	)

	if err := redisdb.Get().SetKey(ctx, key, value, timeout); err != nil {
		t.Error(err)
		return
	}
}

func TestGetValue(t *testing.T) {
	TestSetRDB(t)
	var (
		ctx = context.Background()
		key = "4176ae56-e924-445c-9407-007eca2cae7f"
	)

	val, err := redisdb.Get().GetValue(ctx, key)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(val)
}
