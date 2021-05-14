package redisdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/thanhpp/prom/cmd/portal/repository/redisdb"
)

func TestSetRDB(t *testing.T) {
	var (
		conf = redisdb.RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}
	)

	if err := redisdb.Set(conf); err != nil {
		t.Error(err)
		return
	}
}

func TestGetValue(t *testing.T) {
	TestSetRDB(t)
	var (
		ctx = context.Background()
		key = "testKey"
	)

	val, err := redisdb.Get().GetValue(ctx, key)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(val)
}
