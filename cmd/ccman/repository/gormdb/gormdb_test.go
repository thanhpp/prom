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
		models = []interface{}{&ccmanrpc.Column{}, &ccmanrpc.Card{}}
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
			Title:       "test1",
			Description: "null",
			ColumnID:    1,
			AssignedTo:  1,
			CreatedBy:   1,
			DueDate:     timerpc.ToTimeRPC(time.Now()),
		}
	)

	if _, err := gormdb.GetGormDB().CreateCard(ctx, card); err != nil {
		t.Error(err)
		return
	}
	fmt.Println(card.ID)
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

func TestCreateColumn(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx    = context.Background()
		column = &ccmanrpc.Column{
			ProjectID: 1,
			Title:     "test2",
			Index:     "",
		}
	)

	if _, err := gormdb.GetGormDB().CreateColumn(ctx, column); err != nil {
		t.Error(err)
		return
	}
}

func TestGetAllFromProjectID(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx              = context.Background()
		projectID uint32 = 1
	)

	cols, err := gormdb.GetGormDB().GetAllFromProjectID(ctx, projectID)
	if err != nil {
		t.Error(err)
		return
	}

	for i := range cols {
		for j := range cols[i].Cards {
			fmt.Println(cols[i].Cards[j])
		}
		fmt.Println()
	}
}

func TestMoveCardToCol(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx             = context.Background()
		cardID   uint32 = 2
		newColID uint32 = 2
	)

	if err := gormdb.GetGormDB().MoveCardToCol(ctx, cardID, newColID); err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteColumnByIDAndMove(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx             = context.Background()
		colID    uint32 = 2
		newColID uint32 = 1
	)

	if err := gormdb.GetGormDB().DeleteColumnByIDAndMove(ctx, colID, newColID); err != nil {
		t.Error(err)
		return
	}
}
