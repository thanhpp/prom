package gormdb_test

import (
	"context"
	"testing"

	"github.com/thanhpp/prom/cmd/usrman/repository/entity"
	"github.com/thanhpp/prom/cmd/usrman/repository/gormdb"
)

func TestInitConnection(t *testing.T) {
	var (
		dsn = "user=thanhpp password=testthanhpp dbname=promuser host=127.0.0.1 port=5432 sslmode=disable"
	)

	if err := gormdb.GetGormDB().InitDBConnection(dsn, "INFO"); err != nil {
		t.Error(err)
		return
	}
}

func TestAutoMigrate(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx    = context.Background()
		models = []interface{}{entity.User{}, entity.Project{}, entity.Team{}}
	)

	if err := gormdb.GetGormDB().AutoMigrate(ctx, models...); err != nil {
		t.Error(err)
		return
	}
}
