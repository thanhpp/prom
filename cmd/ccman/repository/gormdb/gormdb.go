package gormdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/thanhpp/prom/cmd/ccman/repository/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- GORMDB ----------------------------------------------------------

type implGorm struct{}

var (
	gDB         = &gorm.DB{}
	gormObj     = new(implGorm)
	cardModel   = new(entity.Card)
	columnModel = new(entity.Column)
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
func (g *implGorm) CreateCard(ctx context.Context, card *entity.Card) (err error) {
	if err := gDB.Table("card").WithContext(ctx).Save(card).Error; err != nil {
		return err
	}

	return nil
}

// GetCardByID ...
func (g *implGorm) GetCardByID(ctx context.Context, cardID uint64) (card *entity.Card, err error) {
	card = new(entity.Card)

	if err := gDB.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Take(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

// GetCardsByDueDate ...
func (g *implGorm) GetCardsByDueDate(ctx context.Context, dueDate time.Time) (cards []*entity.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("due_date = ?", dueDate).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByAssignedToID ...
func (g *implGorm) GetCardsByAssignedToID(ctx context.Context, assignedTo uint64) (cards []*entity.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("assigned_to = ?", assignedTo).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByCreatorID ...
func (g *implGorm) GetCardsByCreatorID(ctx context.Context, creatorID uint64) (cards []*entity.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("created_by = ?", creatorID).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardsByColumnID ...
func (g *implGorm) GetCardsByColumnID(ctx context.Context, columnID uint64) (cards []*entity.Card, err error) {
	rows, err := gDB.Model(cardModel).WithContext(ctx).Where("column_id = ?", columnID).Rows()
	if err != nil {
		return nil, err
	}

	cards, err = scanCards(rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// UpdateCardByID ...
func (g *implGorm) UpdateCardByID(ctx context.Context, cardID uint64, card *entity.Card) (err error) {
	if err = gDB.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Updates(card).Error; err != nil {
		return err
	}

	return nil
}

// DeleteCardByID ...
func (g *implGorm) DeleteCardByID(ctx context.Context, cardID uint64) (err error) {
	if err = gDB.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Delete(&entity.Card{}).Error; err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- COLUMN ----------------------------------------------------------

// CreateColumn ...
func (g *implGorm) CreateColumn(ctx context.Context, column *entity.Column) (err error) {
	if err = gDB.Model(columnModel).WithContext(ctx).Save(column).Error; err != nil {
		return err
	}

	return nil
}

// GetColumnByID ...
func (g *implGorm) GetColumnByID(ctx context.Context, columnID uint64) (column *entity.Column, err error) {
	column = new(entity.Column) //alloc

	if err = gDB.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Take(column).Error; err != nil {
		return nil, err
	}

	return column, nil
}

// GetColumnsByTitle ...
func (g *implGorm) GetColumnsByTitle(ctx context.Context, title string) (columns []*entity.Column, err error) {
	rows, err := gDB.Model(columnModel).WithContext(ctx).Where("title LIKE %?%", fmt.Sprintf("%%%s%%", title)).Rows()
	if err != nil {
		return nil, err
	}

	columns, err = scanColumns(rows)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// GetColumnsByProjectID ...
func (g *implGorm) GetColumnsByProjectID(ctx context.Context, projectID uint64) (columns []*entity.Column, err error) {
	rows, err := gDB.Model(columnModel).WithContext(ctx).Where("project_id = ?", projectID).Rows()
	if err != nil {
		return nil, err
	}

	columns, err = scanColumns(rows)
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// UpdateColumnByID ...
func (g *implGorm) UpdateColumnByID(ctx context.Context, columnID uint64, column *entity.Column) (err error) {
	if err = gDB.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Updates(column).Error; err != nil {
		return err
	}

	return nil
}

// DeleteColumnByID ...
func (g *implGorm) DeleteColumnByID(ctx context.Context, columnID uint64) (err error) {
	if err = gDB.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Delete(&entity.Column{}).Error; err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UTILS ----------------------------------------------------------

func scanCards(rows *sql.Rows) (cards []*entity.Card, err error) {
	for rows.Next() {
		var card = new(entity.Card)
		if err := rows.Scan(card); err != nil {
			return nil, err
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func scanColumns(rows *sql.Rows) (columns []*entity.Column, err error) {
	for rows.Next() {
		var column = new(entity.Column)
		if err := rows.Scan(column); err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil
}
