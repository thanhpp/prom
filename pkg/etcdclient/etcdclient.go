package etcdclient

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDClient struct {
	client *clientv3.Client
}

type ETCDConfigs struct {
	Endpoints []string `mapstructure:"endpoints"`
	Timeout   int      `mapstructure:"timeout"`
}

var ecli = new(ETCDClient)

func Set(econf *ETCDConfigs) (err error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   econf.Endpoints,
		DialTimeout: time.Second * time.Duration(econf.Timeout),
	})
	if err != nil {
		return err
	}

	ecli.client = cli

	return nil
}

func Get() *ETCDClient {
	return ecli
}

func (ec *ETCDClient) SaveKeyValue(ctx context.Context, key string, val string) (err error) {
	_, err = ec.client.Put(ctx, key, val)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ETCDClient) GetValueByKey(ctx context.Context, key string) (val []string, err error) {
	resp, err := ec.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	for i := range resp.Kvs {
		val = append(val, string(resp.Kvs[i].Value))
	}

	return val, nil
}

func (ec *ETCDClient) Close() (err error) {
	if err = ec.client.Close(); err != nil {
		return err
	}

	return nil
}
