package router_test

import (
	"testing"

	"github.com/thanhpp/prom/cmd/portal/webserver/router"
)

func TestNewRouter(t *testing.T) {
	r := router.NewRouter()
	if err := r.Run(":12345"); err != nil {
		t.Error(err)
		return
	}
}
