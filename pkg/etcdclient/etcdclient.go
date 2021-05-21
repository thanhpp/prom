package etcdclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"

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
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
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

func (ec *ETCDClient) RemoveEndpoints(ctx context.Context, service string, address string) (err error) {
	r := &etcdnaming.GRPCResolver{Client: ec.client}
	if err := r.Update(ctx, service, naming.Update{Op: naming.Delete, Addr: address}); err != nil {
		return err
	}

	return nil
}

func (ec *ETCDClient) AddEndpoints(ctx context.Context, service string, address string) (err error) {
	r := &etcdnaming.GRPCResolver{Client: ec.client}
	fmt.Printf("Srv: %s. Addr: %s", service, address)
	if err := r.Update(ctx, service, naming.Update{Op: naming.Delete, Addr: address}); err != nil {
		return err
	}
	if err := r.Update(ctx, service, naming.Update{Op: naming.Add, Addr: address}); err != nil {
		return err
	}

	return nil
}

type Service struct {
	Name string `json:"-"`
	Addr string `json:"addr"`
}

func (ec *ETCDClient) GetServices(ctx context.Context, servicePrefix string) (services []*Service, err error) {
	resp, err := ec.client.Get(ctx, servicePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for i := range resp.Kvs {
		srv := new(Service)
		if err = json.Unmarshal(resp.Kvs[i].Value, srv); err != nil {
			return nil, err
		}
		srv.Name = strings.Split(string(resp.Kvs[i].Key), "/")[0]
		services = append(services, srv)
	}

	return services, nil
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
