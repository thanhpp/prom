package redisdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// ----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- REDIS DB ----------------------------------------------------------

type rDB struct {
	client *redis.Client
}

type iRedisDB interface {
	Set(conf RedisConfig) (err error)
	SetKey(ctx context.Context, key string, value string, timeout time.Duration) (err error)
	GetValue(ctx context.Context, key string) (value string, err error)
	DeleteKey(ctx context.Context, key string) (err error)
}

var iplmRDB = new(rDB)

func Get() iRedisDB {
	return iplmRDB
}

type RedisConfig struct {
	Addr     string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (r *rDB) Set(conf RedisConfig) (err error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	iplmRDB.client = cli

	return nil
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNCTIONS ----------------------------------------------------------

func (r *rDB) SetKey(ctx context.Context, key string, value string, timeout time.Duration) (err error) {
	if err := r.client.Set(ctx, key, value, timeout).Err(); err != nil {
		return nil
	}
	return nil
}

func (r *rDB) GetValue(ctx context.Context, key string) (value string, err error) {
	value, err = r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *rDB) DeleteKey(ctx context.Context, key string) (err error) {
	_, err = r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
