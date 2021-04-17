package repository

import (
	"context"
	"time"

	"github.com/thanhpp/prom/cmd/ccman/repository/entity"
)

type iDao interface {
	InitDBConnection(dsn string, logLevel string) (err error)
	AutoMigrate(ctx context.Context, models ...interface{}) (err error)

	GetCardByID(ctx context.Context, cardID uint64) (card *entity.Card, err error)
	GetCardByDueDate(ctx context.Context, dueDate time.Time) (cards []*entity.Card, err error)
	UpdateCardByID(ctx context.Context, cardID uint64, card *entity.Card) (err error)
	DeleteCardByID(ctx context.Context, cardID uint64) (err error)
}
