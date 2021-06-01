package gormdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	gDB           = &gorm.DB{}
	gormObj       = new(implGorm)
	cardModel     = new(ccmanrpc.Card)
	columnModel   = new(ccmanrpc.Column)
	ErrSameColumn = errors.New("Same column ID")
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
func (g *implGorm) AutoMigrate(ctx context.Context) (err error) {
	if err = g.autoMigrate(ctx, &ccmanrpc.Column{}, &ccmanrpc.Card{}); err != nil {
		return err
	}

	return nil
}
func (g *implGorm) autoMigrate(ctx context.Context, models ...interface{}) (err error) {
	if err = gDB.WithContext(ctx).AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CARD ----------------------------------------------------------

// GetAllFromProjectID

func (g *implGorm) GetAllFromProjectID(ctx context.Context, projectID uint32) (cols []*ccmanrpc.Column, err error) {
	if err := gDB.WithContext(ctx).Model(columnModel).
		Preload("Cards", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("card.index ASC")
		}).
		Where("project_id = ?", projectID).
		Order("project_index ASC").
		Find(&cols).Error; err != nil {
		return nil, err
	}

	// for i := range cols {
	// 	index := strings.Split(strings.TrimRight(cols[i].Index, ","), ",")
	// 	if len(cols[i].Cards) == 0 {
	// 		continue
	// 	}
	// 	if len(index) != len(cols[i].Cards) {
	// 		return nil, errors.New("Mismatch index and cards len")
	// 	}

	// 	for j := range index {
	// 		for k := range cols[i].Cards {
	// 			id, err := strconv.Atoi(index[j])
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			if cols[i].Cards[k].ID == uint32(id) {
	// 				cols[i].Cards = append(cols[i].Cards, cols[i].Cards[k])
	// 			}
	// 		}
	// 	}

	// 	cols[i].Cards = cols[i].Cards[len(cols[i].Cards)/2:]
	// }

	for i := range cols {
		fmt.Println(cols[i].ID, cols[i].ProjectIndex)
	}
	return cols, nil
}

// // CreateCard Const
// const createCardColIndex = "UPDATE \"column\" SET index = ? || ',' || index  WHERE id = ?"

// CreateCard ...
func (g *implGorm) CreateCard(ctx context.Context, card *ccmanrpc.Card) (createdID uint32, err error) {
	err = gDB.Transaction(func(tx *gorm.DB) error {
		// if err := tx.WithContext(ctx).Table("card").WithContext(ctx).Save(card).Error; err != nil {
		// 	return err
		// }

		// if err := tx.WithContext(ctx).Exec(createCardColIndex, strconv.Itoa(int(card.ID)), card.ColumnID).Error; err != nil {
		// 	return err
		// }

		// return nil
		_, err = g.createCardv2(ctx, tx, card)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return card.ID, nil
}

func (g *implGorm) createCardv2(ctx context.Context, tx *gorm.DB, card *ccmanrpc.Card) (cardID uint32, err error) {
	// get column max index & increase
	var maxIndex uint32
	if err = tx.WithContext(ctx).Model(columnModel).Where("id = ?", card.ColumnID).Pluck("max_index", &maxIndex).Limit(1).Error; err != nil {
		return 0, err
	}
	maxIndex = maxIndex + 1
	if err = tx.WithContext(ctx).Model(columnModel).Where("id = ?", card.ColumnID).Update("max_index", maxIndex).Error; err != nil {
		return 0, err
	}

	// save card
	card.Index = maxIndex
	if err = tx.WithContext(ctx).Model(cardModel).Save(card).Error; err != nil {
		return 0, err
	}

	return card.ID, nil
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

// MoveCardToCol
func (g *implGorm) MoveCardToCol(ctx context.Context, cardID uint32, newColID uint32) (err error) {
	err = gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := moveCardToColTransaction(tx, ctx, cardID, newColID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ----------
// MoveCardToColv2
func (g *implGorm) MoveCardToColv2(ctx context.Context, cardID uint32, colID uint32, aboveIdx uint32) (err error) {
	// pre-exec check
	if aboveIdx == 0 {
		return errors.New("Zero index")
	}

	err = gDB.Transaction(func(tx *gorm.DB) error {
		if err := g.moveCardToColv2(ctx, tx, cardID, colID, aboveIdx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *implGorm) moveCardToColv2(ctx context.Context, tx *gorm.DB, cardID uint32, colID uint32, aboveIdx uint32) (err error) {
	// precheck max_index
	var maxIndex uint32
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("id = ?", colID).
		Pluck("max_index", &maxIndex).Limit(1).Error; err != nil {
		return err
	}
	if aboveIdx > maxIndex+1 || aboveIdx < 1 {
		return errors.New("Invalid new index")
	}

	// remove from old column, decrease all larger(newer) index
	var card = new(ccmanrpc.Card)
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("id = ?", cardID).
		Take(card).Error; err != nil {
		return err
	}

	if card.ColumnID == colID {
		return ErrSameColumn
	}

	if err = tx.WithContext(ctx).Model(cardModel).
		Where("column_id = ? AND index > ?", card.ColumnID, card.Index).
		Update("index", gorm.Expr("index - ?", 1)).Error; err != nil {
		return err
	}
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("id = ?", card.ColumnID).
		Update("max_index", gorm.Expr("max_index - ?", 1)).Error; err != nil {
		return err
	}

	// add to new column, create empty index slot & update card
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("column_id = ? AND index >= ?", colID, aboveIdx).
		Update("index", gorm.Expr("index + ?", 1)).Error; err != nil {
		return err
	}
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("id = ?", colID).
		Update("max_index", gorm.Expr("max_index + ?", 1)).Error; err != nil {
		return err
	}
	card.ColumnID = colID
	card.Index = aboveIdx
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("id = ?", card.ID).
		Updates(card).Error; err != nil {
		return err
	}
	return nil
}

// UpdateCardByID ..
func (g *implGorm) UpdateCardByID(ctx context.Context, cardID uint32, card *ccmanrpc.Card) (err error) {
	err = gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// if card.ColumnID > 0 {
		// 	dbCard := new(ccmanrpc.Card)
		// 	if err = tx.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Take(dbCard).Error; err != nil {
		// 		return err
		// 	}

		// 	if card.ColumnID != dbCard.ID {
		// 		if err = moveCardToColTransaction(tx, ctx, cardID, card.ColumnID); err != nil {
		// 			return err
		// 		}
		// 	}
		// }
		card.ColumnID = 0 // prevent wrong column update
		card.Index = 0    // prevent wrong index update
		if err = tx.Model(cardModel).WithContext(ctx).
			Where("id = ?", cardID).
			Updates(card).Error; err != nil {
			return err
		}

		return nil
	})

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
func (g *implGorm) CreateColumn(ctx context.Context, column *ccmanrpc.Column) (createdID uint32, err error) {
	err = gDB.Transaction(func(tx *gorm.DB) error {
		maxIdx, er := g.maxPrjColIndex(ctx, tx, column.ProjectID)
		if er != nil {
			return er
		}
		column.ProjectIndex = maxIdx + 1

		if er = tx.Model(columnModel).WithContext(ctx).Save(column).Error; err != nil {
			return er
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return column.ID, nil
}

func (g *implGorm) maxPrjColIndex(ctx context.Context, tx *gorm.DB, prjID uint32) (maxPrjColIdx uint32, err error) {
	if err = tx.Model(columnModel).
		Select("MAX(project_index) AS max_prj_col_idx").
		Where("project_id = ?", prjID).
		Pluck("max_prj_col_idx", &maxPrjColIdx).Error; err != nil {
		return 0, nil
	}

	return maxPrjColIdx, nil
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
func (g *implGorm) DeleteColumnAndAllCardByID(ctx context.Context, columnID uint32) (err error) {
	err = gDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Model(cardModel).
			Where("column_id = ?", columnID).
			Delete(&ccmanrpc.Column{}).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(columnModel).
			Where("id = ?", columnID).
			Delete(&ccmanrpc.Column{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// DeleteColumnByID ...
func (g *implGorm) DeleteColumnByIDAndMove(ctx context.Context, columnID uint32, newColID uint32) (err error) {
	if columnID == newColID {
		return errors.New("Duplicate columnID")
	}

	err = gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		column := new(ccmanrpc.Column)
		if err := tx.WithContext(ctx).Model(columnModel).Where("id = ?", columnID).Take(column).Error; err != nil {
			return err
		}

		cardIDs := strings.Split(strings.TrimRight(column.Index, ","), ",")
		for i := range cardIDs {
			id, err := strconv.Atoi(cardIDs[i])
			if err != nil {
				return err
			}
			if err := moveCardToColTransaction(tx, ctx, uint32(id), newColID); err != nil {
				return err
			}
		}

		if err = tx.Model(columnModel).WithContext(ctx).Where("id = ?", columnID).Delete(&ccmanrpc.Column{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ------------------------------
// ReorderCard ...
func (g *implGorm) ReorderCard(ctx context.Context, cardID uint32, aboveIdx uint32) (err error) {
	// pre-exec check
	if aboveIdx == 0 {
		return errors.New("Zero index")
	}

	err = gDB.Transaction(func(tx *gorm.DB) error {
		var card = new(ccmanrpc.Card)
		if err := tx.WithContext(ctx).Model(cardModel).
			Where("id = ?", cardID).Take(card).Error; err != nil {
			return err
		}

		if card.Index == aboveIdx {
			return nil
		}

		if aboveIdx > card.Index {
			if err := g.moveCardDown(ctx, tx, card, aboveIdx); err != nil {
				return err
			}
		} else {
			if err := g.moveCardUp(ctx, tx, card, aboveIdx); err != nil {
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

func (g *implGorm) moveCardDown(ctx context.Context, tx *gorm.DB, card *ccmanrpc.Card, aboveIdx uint32) (err error) {
	// pre-exec check
	if card.Index >= aboveIdx {
		return errors.New("Invalid index")
	}
	var maxIndex uint32
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("id = ?", card.ColumnID).
		Pluck("max_index", &maxIndex).
		Limit(1).Error; err != nil {
		return err
	}
	if aboveIdx > maxIndex+1 {
		return errors.New("Invalid index")
	}

	// create empty slot
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("column_id = ? AND index > ? AND index < ?", card.ColumnID, card.Index, aboveIdx).
		Update("index", gorm.Expr("index - ?", 1)).Error; err != nil {
		return err
	}

	// update card new index
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("id = ?", card.ID).
		Update("index", aboveIdx-1).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) moveCardUp(ctx context.Context, tx *gorm.DB, card *ccmanrpc.Card, aboveIdx uint32) (err error) {
	// pre-exec check
	if card.Index <= aboveIdx || aboveIdx < 1 {
		return errors.New("Invalid index")
	}

	// create empty slot
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("column_id = ? AND index >= ? AND index < ?", card.ColumnID, aboveIdx, card.Index).
		Update("index", gorm.Expr("index + ?", 1)).Error; err != nil {
		return err
	}

	// update card index
	if err = tx.WithContext(ctx).Model(cardModel).
		Where("id = ?", card.ID).
		Update("index", aboveIdx).Error; err != nil {
		return err
	}

	return nil
}

// ------------------------------
// Reorder column
func (g *implGorm) ReorderColumn(ctx context.Context, columnID uint32, nextOfIdx uint32) (err error) {
	err = gDB.Transaction(func(tx *gorm.DB) error {
		var column = new(ccmanrpc.Column)
		if err := tx.WithContext(ctx).Model(columnModel).
			Where("id = ?", columnID).Take(column).Error; err != nil {
			return err
		}

		if column.ProjectIndex == nextOfIdx {
			return nil
		}

		if column.ProjectIndex < nextOfIdx {
			if err := g.moveColumnForward(ctx, tx, column, nextOfIdx); err != nil {
				return err
			}
		} else {
			if err := g.moveColumnBackward(ctx, tx, column, nextOfIdx); err != nil {
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

func (g *implGorm) moveColumnForward(ctx context.Context, tx *gorm.DB, column *ccmanrpc.Column, nextOfIdx uint32) (err error) {
	// pre-exec check
	if column.ProjectIndex >= nextOfIdx {
		return errors.New("Invalid Index")
	}

	maxColPrjIdx, err := g.maxPrjColIndex(ctx, tx, column.ProjectID)
	if err != nil {
		return err
	}
	if nextOfIdx > maxColPrjIdx {
		return errors.New("Invalid Index")
	}

	// create slot
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("project_id = ? AND project_index > ? AND project_index <= ?", column.ProjectID, column.ProjectIndex, nextOfIdx).
		Update("project_index", gorm.Expr("project_index - ?", 1)).Error; err != nil {
		return err
	}

	// update column index
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("id = ?", column.ID).
		Update("project_index", nextOfIdx).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) moveColumnBackward(ctx context.Context, tx *gorm.DB, column *ccmanrpc.Column, nextOfIdx uint32) (err error) {
	// pre-exec check
	if column.ProjectIndex <= nextOfIdx {
		return errors.New("Invalid Index")
	}

	// create slot
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("project_id = ? AND project_index < ? AND project_index > ?", column.ProjectID, column.ProjectIndex, nextOfIdx).
		Update("project_index", gorm.Expr("project_index + ?", 1)).Error; err != nil {
		return err
	}

	// update column index
	if err = tx.WithContext(ctx).Model(columnModel).
		Where("id = ?", column.ID).
		Update("project_index", nextOfIdx+1).Error; err != nil {
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

// moveCardToColTransaction
const moveCardColIndex = "UPDATE \"column\" SET index = ? || ',' || index  WHERE id = ?"

func moveCardToColTransaction(tx *gorm.DB, ctx context.Context, cardID uint32, newColID uint32) (err error) {
	card := new(ccmanrpc.Card)
	if err := tx.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Take(card).Error; err != nil {
		return err
	}

	// remove from old column index
	column := new(ccmanrpc.Column)
	if err := tx.Model(columnModel).WithContext(ctx).Where("id = ?", card.ColumnID).Take(column).Error; err != nil {
		return err
	}

	column.Index = strings.Replace(column.Index, fmt.Sprintf("%d,", card.ID), "", 1)
	if len(column.Index) > 0 && column.Index[len(column.Index)-1] != ',' {
		column.Index += ","
	}

	if err := tx.Model(columnModel).WithContext(ctx).Where("id = ?", column.ID).Updates(column).Error; err != nil {
		return err
	}

	// add to new column index
	if err := tx.WithContext(ctx).Exec(moveCardColIndex, strconv.Itoa(int(card.ID)), newColID).Error; err != nil {
		return err
	}

	// update card columnID
	if err := tx.Model(cardModel).WithContext(ctx).Where("id = ?", cardID).Update("column_id", newColID).Error; err != nil {
		return err
	}

	return nil
}
