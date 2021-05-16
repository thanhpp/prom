package etcdclient_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/thanhpp/prom/pkg/etcdclient"
)

func Init() {
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
	Init()
	var (
		ctx = context.Background()
		key = "testkey"
		val = "testvalue"
	)

	if err := etcdclient.Get().SaveKeyValue(ctx, key, val); err != nil {
		t.Error(err)
		return
	}
}

func TestGetValueByKey(t *testing.T) {
	Init()
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
