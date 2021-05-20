package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/thanhpp/prom/cmd/portal/service"
)

func TestSetUsrManService(t *testing.T) {
	var (
		ctx    = context.Background()
		target = "127.0.0.1:8090"
	)

	if err := service.SetUsrManService(ctx, target); err != nil {
		t.Error(err)
		return
	}

}

func TestNewUser(t *testing.T) {
	TestSetUsrManService(t)

	var (
		ctx      = context.Background()
		username = "testusername1"
		pass     = "testpass1"
	)

	if err := service.GetUsrManService().NewUser(ctx, username, pass); err != nil {
		t.Error(err)
		return
	}
}

func TestLogin(t *testing.T) {
	TestSetUsrManService(t)

	var (
		ctx      = context.Background()
		username = "testusername1"
		pass     = "testpass1"
	)

	user, err := service.GetUsrManService().Login(ctx, username, pass)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(user)
}
