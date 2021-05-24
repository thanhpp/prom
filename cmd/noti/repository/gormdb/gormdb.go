package gormdb

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/thanhpp/prom/cmd/noti/repository/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type implGorm struct{}

var (
	gDB           = &gorm.DB{}
	gormObj       = new(implGorm)
	notiModel     = &entity.Notification{}
	userNotiModel = &entity.UserNoti{}
)

// GetGormDB ...
func GetGormDB() *implGorm {
	return gormObj
}

// InitDBConnection ...
func (g *implGorm) InitDBConnection(dsn string, logLevel string) (err error) {
	var (
		gormConfig = &gorm.Config{
			Logger: gormlog.Default.LogMode(gormlog.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
	)

	// create db log
	switch logLevel {
	case "INFO":
		gormConfig.Logger = gormlog.Default.LogMode(gormlog.Info)
	case "WARN":
		gormConfig.Logger = gormlog.Default.LogMode(gormlog.Warn)
	case "ERROR":
		gormConfig.Logger = gormlog.Default.LogMode(gormlog.Error)
	case "SILENT":
		gormConfig.Logger = gormlog.Default.LogMode(gormlog.Silent)
	default:
		log.Println("START GORM LOG WITH DEFAULT CONFIG: INFO")
	}

	gDB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return err
	}

	return nil
}

// AutoMigrate ...
func (g *implGorm) AutoMigrate(ctx context.Context, models ...interface{}) (err error) {
	if err = gDB.WithContext(ctx).AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- NOTIFICATIONS ----------------------------------------------------------

func (g *implGorm) CreateNotification(ctx context.Context, noti *entity.Notification, users []int) (err error) {
	err = gDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(notiModel).Save(noti).Error; err != nil {
			return err
		}

		var userNotis = make([]*entity.UserNoti, 0, len(users))
		for i := range users {
			userNotis = append(userNotis, &entity.UserNoti{
				UserID: users[i],
				NotiID: noti.ID,
				Seen:   false,
			})
		}

		if err := tx.WithContext(ctx).Model(userNotiModel).CreateInBatches(userNotis, len(userNotis)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *implGorm) GetNotiByUserID(ctx context.Context, userID int, page int, size int) (notis []*entity.Notification, err error) {
	var (
		offset = page * size
		limit  = size
	)
	if page <= 0 || size <= 0 {
		offset = 1
		limit = 1
	}

	userNotiRows, err := gDB.WithContext(ctx).Model(userNotiModel).Where("user_id = ?", userID).Order("noti_id DESC").Offset(offset).Limit(limit).Rows()
	if err != nil {
		return nil, err
	}
	userNotis, err := scanUserNotis(ctx, gDB, userNotiRows)
	if err != nil {
		return nil, err
	}

	var (
		notiIDs = make([]int, 0, len(userNotis))
		seen    = make([]bool, 0, len(userNotis))
	)
	for i := range userNotis {
		notiIDs = append(notiIDs, userNotis[i].NotiID)
		seen = append(seen, userNotis[i].Seen)
	}

	notiRows, err := gDB.WithContext(ctx).Model(notiModel).Where("id IN ?", notiIDs).Order("id DESC").Rows()
	if err != nil {
		return nil, err
	}

	notis, err = scanNoti(ctx, gDB, notiRows, seen)
	if err != nil {
		return nil, err
	}

	return notis, nil
}

func (g *implGorm) GetNotiByCardID(ctx context.Context, cardID int, page int, size int) (notis []*entity.Notification, err error) {
	var (
		offset = page * size
		limit  = size
	)
	if page <= 0 || size <= 0 {
		offset = 1
		limit = 1
	}

	rows, err := gDB.WithContext(ctx).Model(notiModel).Where("card_id = ?", cardID).Order("id DESC").Offset(offset).Limit(limit).Rows()
	if err != nil {
		return nil, err
	}

	notis, err = scanNoti(ctx, gDB, rows, nil)
	if err != nil {
		return nil, err
	}

	return notis, nil
}

func (g *implGorm) UpdateSeen(ctx context.Context, userID int, notiID int, seen bool) (err error) {
	if err = gDB.WithContext(ctx).Model(userNotiModel).Update("seen", seen).Where("user_id = ? AND noti_id = ?", userID, notiID).Error; err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UTILS ----------------------------------------------------------

func scanUserNotis(ctx context.Context, tx *gorm.DB, rows *sql.Rows) (userNotis []*entity.UserNoti, err error) {
	for rows.Next() {
		var userNoti = new(entity.UserNoti)
		if err = tx.WithContext(ctx).ScanRows(rows, userNoti); err != nil {
			return nil, err
		}
		userNotis = append(userNotis, userNoti)
	}

	return userNotis, nil
}

func scanNoti(ctx context.Context, tx *gorm.DB, rows *sql.Rows, seen []bool) (notis []*entity.Notification, err error) {
	for rows.Next() {
		var noti = new(entity.Notification)
		if err = tx.WithContext(ctx).ScanRows(rows, noti); err != nil {
			return nil, err
		}

		notis = append(notis, noti)
	}

	if seen != nil {
		if len(seen) != len(notis) {
			return nil, errors.New("Mismatch len noti and seen")
		}

		for i := range seen {
			notis[i].Seen = seen[i]
		}
	}

	return notis, nil
}
