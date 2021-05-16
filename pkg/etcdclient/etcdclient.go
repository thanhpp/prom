package etcdclient

import (
	"context"
	"time"

	"google.golang.org/grpc/naming"

	clientv3 "github.com/coreos/etcd/clientv3"
	etcdnaming "go.etcd.io/etcd/clientv3/naming"
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

func (ec *ETCDClient) AddEndpoints(ctx context.Context, service string, address string) (err error) {
	r := &etcdnaming.GRPCResolver{Client: ec.client}
	if err := r.Update(ctx, service, naming.Update{Op: naming.Add, Addr: address}); err != nil {
		return err
	}

	return nil
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

func (ec *ETCDClient) Resolver() *etcdnaming.GRPCResolver {
	return &etcdnaming.GRPCResolver{Client: ec.client}
}
