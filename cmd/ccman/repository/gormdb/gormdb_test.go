package gormdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/thanhpp/prom/cmd/ccman/repository/gormdb"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/timerpc"
)

func TestInitConnection(t *testing.T) {
	var (
		dsn = "user=thanhpp password=testthanhpp dbname=prom host=127.0.0.1 port=5432 sslmode=disable"
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
		models = []interface{}{ccmanrpc.Card{}, ccmanrpc.Column{}}
	)

	if err := gormdb.GetGormDB().AutoMigrate(ctx, models...); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateCard(t *testing.T) {
	TestInitConnection(t)

	var (
		ctx  = context.Background()
		card = &ccmanrpc.Card{
			Title:       "test",
			Description: "null",
			ColumnID:    1,
			AssignedTo:  1,
			CreatedBy:   1,
			DueDate:     timerpc.ToTimeRPC(time.Now()),
		}
	)

	if err := gormdb.GetGormDB().CreateCard(ctx, card); err != nil {
		t.Error(err)
		return
	}
}

func TestGetCardByID(t *testing.T) {
	TestInitConnection(t)

	var (
		ctx        = context.Background()
		id  uint32 = 1
	)

	card, err := gormdb.GetGormDB().GetCardByID(ctx, id)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%+v\n", card)
	fmt.Println(timerpc.ToTime(card.CreatedAt))
}

func TestGetCardByDueDate(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx     = context.Background()
		duedate = timerpc.ToTime(&timerpc.Time{
			Seconds: 1619005776,
			Nanos:   210219000,
		})
	)

	cards, err := gormdb.GetGormDB().GetCardsByDueDate(ctx, duedate)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(cards)
}
