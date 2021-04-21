package repository

import (
	"context"
	"time"

	"github.com/thanhpp/prom/cmd/ccman/repository/entity"
	"github.com/thanhpp/prom/cmd/ccman/repository/gormdb"
)

type iDao interface {
	InitDBConnection(dsn string, logLevel string) (err error)
	AutoMigrate(ctx context.Context, models ...interface{}) (err error)

	CreateCard(ctx context.Context, card *entity.Card) (err error)
	GetCardByID(ctx context.Context, cardID uint64) (card *entity.Card, err error)
	GetCardsByDueDate(ctx context.Context, dueDate time.Time) (cards []*entity.Card, err error)
	GetCardsByAssignedToID(ctx context.Context, assignedTo uint64) (cards []*entity.Card, err error)
	GetCardsByCreatorID(ctx context.Context, creatorID uint64) (cards []*entity.Card, err error)
	GetCardsByColumnID(ctx context.Context, columnID uint64) (cards []*entity.Card, err error)
	UpdateCardByID(ctx context.Context, cardID uint64, card *entity.Card) (err error)
	DeleteCardByID(ctx context.Context, cardID uint64) (err error)

	CreateColumn(ctx context.Context, column *entity.Column) (err error)
	GetColumnByID(ctx context.Context, columnID uint64) (column *entity.Column, err error)
	GetColumnsByTitle(ctx context.Context, title string) (columns []*entity.Column, err error)
	GetColumnsByProjectID(ctx context.Context, projectID uint64) (columns []*entity.Column, err error)
	UpdateColumnByID(ctx context.Context, columnID uint64, column *entity.Column) (err error)
	DeleteColumnByID(ctx context.Context, columnID uint64) (err error)
}

var dao iDao = gormdb.GetGormDB()

func GetDAO() iDao {
	return dao
}
