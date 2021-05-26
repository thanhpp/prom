package repository

import (
	"context"

	"github.com/thanhpp/prom/cmd/noti/repository/entity"
	"github.com/thanhpp/prom/cmd/noti/repository/gormdb"
)

type iDAO interface {
	InitDBConnection(dsn string, logLevel string) (err error)
	AutoMigrate(ctx context.Context, models ...interface{}) (err error)
	CreateNotification(ctx context.Context, noti *entity.Notification, users []int) (err error)
	GetNotiByUserID(ctx context.Context, userID int, page int, size int) (notis []*entity.Notification, err error)
	GetNotiByCardID(ctx context.Context, cardID int, page int, size int) (notis []*entity.Notification, err error)
	UpdateSeen(ctx context.Context, userID int, notiID int, seen bool) (err error)
}

func Get() iDAO {
	return gormdb.GetGormDB()
}
