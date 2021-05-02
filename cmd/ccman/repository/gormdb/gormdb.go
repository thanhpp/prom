package gormdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/thanhpp/prom/pkg/ccmanrpc"
)

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- GORMDB ----------------------------------------------------------

type implGorm struct{}

var (
	gDB         = &gorm.DB{}
	gormObj     = new(implGorm)
	cardModel   = new(ccmanrpc.Card)
	columnModel = new(ccmanrpc.Column)
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
	err = gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range models {
			if err := tx.WithContext(ctx).AutoMigrate(models[i]); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CARD ----------------------------------------------------------

// CreateCard ...
func (g *implGorm) CreateCard(ctx context.Context, card *ccmanrpc.Card) (err error) {
	if err := gDB.Table("card").WithContext(ctx).Save(card).Error; err != nil {
		return err
	}

	return nil
}

// GetCardByID ...
func (g *implGorm) GetCardByID(ctx context.Context, cardID uint32) (card *ccmanrpc.Card, err error) {
	card = new(ccmanrpc.Card)

	if err := gDB.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Take(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

// GetCardsByDueDate ...
func (g *implGorm) GetCardsByDueDate(ctx context.Context, dueDate time.Time) (cards []*ccmanrpc.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("due_date = ?", dueDate).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(gDB, rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByAssignedToID ...
func (g *implGorm) GetCardsByAssignedToID(ctx context.Context, assignedTo uint32) (cards []*ccmanrpc.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("assigned_to = ?", assignedTo).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(gDB, rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByCreatorID ...
func (g *implGorm) GetCardsByCreatorID(ctx context.Context, creatorID uint32) (cards []*ccmanrpc.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("created_by = ?", creatorID).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(gDB, rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByColumnID ...
func (g *implGorm) GetCardsByColumnID(ctx context.Context, columnID uint32) (cards []*ccmanrpc.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("column_id = ?", columnID).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(gDB, rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// UpdateCardByID ...
func (g *implGorm) UpdateCardByID(ctx context.Context, cardID uint32, card *ccmanrpc.Card) (err error) {
	if err = gDB.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Updates(card).Error; err != nil {
		return err
	}

	return nil
}

// DeleteCardByID ...
func (g *implGorm) DeleteCardByID(ctx context.Context, cardID uint32) (err error) {
	if err = gDB.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Delete(&ccmanrpc.Card{}).Error; err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- COLUMN ----------------------------------------------------------

// CreateColumn ...
func (g *implGorm) CreateColumn(ctx context.Context, column *ccmanrpc.Column) (err error) {
	if err = gDB.Model(columnModel).WithContext(ctx).Save(column).Error; err != nil {
		return err
	}

	return nil
}

// GetColumnByID ...
func (g *implGorm) GetColumnByID(ctx context.Context, columnID uint32) (column *ccmanrpc.Column, err error) {
	column = new(ccmanrpc.Column) //alloc

	if err = gDB.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Take(column).Error; err != nil {
		return nil, err
	}

	return column, nil
}

// GetColumnsByTitle ...
func (g *implGorm) GetColumnsByTitle(ctx context.Context, title string) (columns []*ccmanrpc.Column, err error) {
	rows, err := gDB.Model(columnModel).WithContext(ctx).Where("title LIKE %?%", fmt.Sprintf("%%%s%%", title)).Rows()
	if err != nil {
		return nil, err
	}

	columns, err = scanColumns(gDB, rows)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// GetColumnsByProjectID ...
func (g *implGorm) GetColumnsByProjectID(ctx context.Context, projectID uint32) (columns []*ccmanrpc.Column, err error) {
	rows, err := gDB.Model(columnModel).WithContext(ctx).Where("project_id = ?", projectID).Rows()
	if err != nil {
		return nil, err
	}

	columns, err = scanColumns(gDB, rows)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// UpdateColumnByID ...
func (g *implGorm) UpdateColumnByID(ctx context.Context, columnID uint32, column *ccmanrpc.Column) (err error) {
	if err = gDB.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Updates(column).Error; err != nil {
		return err
	}

	return nil
}

// DeleteColumnByID ...
func (g *implGorm) DeleteColumnByID(ctx context.Context, columnID uint32) (err error) {
	if err = gDB.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Delete(&ccmanrpc.Column{}).Error; err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UTILS ----------------------------------------------------------

func scanCards(gormDB *gorm.DB, rows *sql.Rows) (cards []*ccmanrpc.Card, err error) {
	for rows.Next() {
		var card = new(ccmanrpc.Card)
		if err := gormDB.ScanRows(rows, card); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func scanColumns(gormDB *gorm.DB, rows *sql.Rows) (columns []*ccmanrpc.Column, err error) {
	for rows.Next() {
		var column = new(ccmanrpc.Column)
		if err := gormDB.ScanRows(rows, column); err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil
}
