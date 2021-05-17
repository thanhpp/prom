package etcdclient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/thanhpp/prom/pkg/etcdclient"
)

func TestInit(t *testing.T) {
	var (
		etcdConfig = &etcdclient.ETCDConfigs{
			Endpoints: []string{"127.0.0.1:2379"},
			Timeout:   5,
		}
	)

	if err := etcdclient.Set(etcdConfig); err != nil {
		panic(err)
	}
}

func TestSaveKeyValue(t *testing.T) {
	TestInit(t)
	var (
		ctx = context.Background()
		key = "testkey"
		val = "testvalue2"
	)

	if err := etcdclient.Get().SaveKeyValue(ctx, key, val); err != nil {
		t.Error(err)
		return
	}
}

func TestGetValueByKey(t *testing.T) {
	TestInit(t)
	var (
		ctx = context.Background()
		key = "testkey"
	)

	val, err := etcdclient.Get().GetValueByKey(ctx, key)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(val)
}

func TestGetServices(t *testing.T) {
	TestInit(t)
	var (
		ctx           = context.Background()
		servicePrefix = "cardscolumnsmanager"
	)

	services, err := etcdclient.Get().GetServices(ctx, servicePrefix)
	if err != nil {
		t.Error(err)
		return
	}

	for i := range services {
		fmt.Println(services[i])
	}
}

func TestRemoveEndpoints(t *testing.T) {
	TestInit(t)
	var (
		ctx     = context.Background()
		service = "cardscolumnsmanager-1"
		addr    = "0.0.0.0:8080"
	)

	if err := etcdclient.Get().RemoveEndpoints(ctx, service, addr); err != nil {
		t.Error(err)
		return
	}
}
