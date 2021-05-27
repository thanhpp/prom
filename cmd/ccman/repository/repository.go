package repository

import (
	"context"
	"time"

	"github.com/thanhpp/prom/cmd/ccman/repository/gormdb"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
)

type iDao interface {
	InitDBConnection(dsn string, logLevel string) (err error)
	AutoMigrate(ctx context.Context) (err error)

	GetAllFromProjectID(ctx context.Context, projectID uint32) (cols []*ccmanrpc.Column, err error)

	CreateCard(ctx context.Context, card *ccmanrpc.Card) (createdID uint32, err error)
	GetCardByID(ctx context.Context, cardID uint32) (card *ccmanrpc.Card, err error)
	GetCardsByDueDate(ctx context.Context, dueDate time.Time) (cards []*ccmanrpc.Card, err error)
	GetCardsByAssignedToID(ctx context.Context, assignedTo uint32) (cards []*ccmanrpc.Card, err error)
	GetCardsByCreatorID(ctx context.Context, creatorID uint32) (cards []*ccmanrpc.Card, err error)
	GetCardsByColumnID(ctx context.Context, columnID uint32) (cards []*ccmanrpc.Card, err error)
	UpdateCardByID(ctx context.Context, cardID uint32, card *ccmanrpc.Card) (err error)
	MoveCardToCol(ctx context.Context, cardID uint32, newColID uint32) (err error)
	MoveCardToColv2(ctx context.Context, cardID uint32, colID uint32, index uint32) (err error)
	DeleteCardByID(ctx context.Context, cardID uint32) (err error)

	CreateColumn(ctx context.Context, column *ccmanrpc.Column) (createdID uint32, err error)
	GetColumnByID(ctx context.Context, columnID uint32) (column *ccmanrpc.Column, err error)
	GetColumnsByTitle(ctx context.Context, title string) (columns []*ccmanrpc.Column, err error)
	GetColumnsByProjectID(ctx context.Context, projectID uint32) (columns []*ccmanrpc.Column, err error)
	UpdateColumnByID(ctx context.Context, columnID uint32, column *ccmanrpc.Column) (err error)
	ReorderCard(ctx context.Context, cardID uint32, newIdx uint32) (err error)
	DeleteColumnByID(ctx context.Context, columnID uint32) (err error)
	DeleteColumnByIDAndMove(ctx context.Context, columnID uint32, newColID uint32) (err error)
}

var dao iDao = gormdb.GetGormDB()

func GetDAO() iDao {
	return dao
}
