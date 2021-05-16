package service

import (
	"context"
	"time"

	"github.com/thanhpp/prom/pkg/timerpc"

	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/errconst"
	"google.golang.org/grpc"
)

// --------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- CCMANService ----------------------------------------------------------

type CCManSrv struct {
	grpcClient ccmanrpc.CCManagerClient
}

func (cs CCManSrv) client() (client ccmanrpc.CCManagerClient) {
	return cs.grpcClient
}

func (cs CCManSrv) name() string {
	return "Card Column manager service"
}

func (cs CCManSrv) respError(err error, respCode int32, respMessage string) error {
	return errconst.ServiceError{Srv: cs.name(), Err: err, Code: respCode, Msg: respMessage}
}

func (cs CCManSrv) error(err error) error {
	return errconst.ServiceError{Srv: cs.name(), Err: err, Msg: err.Error()}
}

type iCCMan interface {
	CreateCard(ctx context.Context, card *ccmanrpc.Card) (err error)
	GetCardByID(ctx context.Context, cardID uint32) (card *ccmanrpc.Card, err error)
	GetCardsByDueDate(ctx context.Context, duedate time.Time) (cards []*ccmanrpc.Card, err error)
	GetCardsByAssignedToID(ctx context.Context, userID uint32) (cards []*ccmanrpc.Card, err error)
	GetCardsByCreatorID(ctx context.Context, userID uint32) (cards []*ccmanrpc.Card, err error)
	GetCardsByColumnID(ctx context.Context, colID uint32) (cards []*ccmanrpc.Card, err error)
	UpdateCardByID(ctx context.Context, cardID uint32, card *ccmanrpc.Card) (err error)
	DeleteCardByID(ctx context.Context, cardID uint32) (err error)
	CreateColumn(ctx context.Context, column *ccmanrpc.Column) (err error)
	GetColumnByID(ctx context.Context, colID uint32) (col *ccmanrpc.Column, err error)
	GetColumnsByTitle(ctx context.Context, title string) (cols []*ccmanrpc.Column, err error)
	GetColumnsByProjectID(ctx context.Context, projectID uint32) (cols []*ccmanrpc.Column, err error)
	UpdateColumnByID(ctx context.Context, colID uint32, col *ccmanrpc.Column) (err error)
	DeleteColumnByID(ctx context.Context, colID uint32) (err error)
}

var implCCManSrv = new(CCManSrv)

func SetCCManSrv(ctx context.Context, target string) (err error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	conn, err := grpc.DialContext(newCtx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		cancel()
		return err
	}

	implCCManSrv.grpcClient = ccmanrpc.NewCCManagerClient(conn)
	return nil
}

func GetCCManSrv() iCCMan {
	return implCCManSrv
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNCTIONS ----------------------------------------------------------

func (cS CCManSrv) CreateCard(ctx context.Context, card *ccmanrpc.Card) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.CreateCardReq{
		CreateCard: card,
	}

	_, err = cS.client().CreateCard(ctx, in)
	if err != nil {
		return cS.error(err)
	}

	return nil
}

func (cS CCManSrv) GetCardByID(ctx context.Context, cardID uint32) (card *ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetCardByIDReq{
		CardID: cardID,
	}

	resp, err := cS.client().GetCardByID(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.ResCard, nil
}

func (cS CCManSrv) GetCardsByDueDate(ctx context.Context, duedate time.Time) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(err)
	}

	in := &ccmanrpc.GetCardsByDueDateReq{
		DueDate: timerpc.ToTimeRPC(duedate),
	}

	resp, err := cS.client().GetCardsByDueDate(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS CCManSrv) GetCardsByAssignedToID(ctx context.Context, userID uint32) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return
	}

	in := &ccmanrpc.GetCardsByAssignedToIDReq{
		AssignedToID: userID,
	}

	resp, err := cS.client().GetCardsByAssignedToID(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS CCManSrv) GetCardsByCreatorID(ctx context.Context, userID uint32) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetCardsByCreatorIDReq{
		CreatorID: userID,
	}

	resp, err := cS.client().GetCardsByCreatorID(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS CCManSrv) GetCardsByColumnID(ctx context.Context, colID uint32) (cards []*ccmanrpc.Card, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetCardsByColumnIDReq{
		ColumnID: colID,
	}

	resp, err := cS.client().GetCardsByColumnID(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Cards, nil
}

func (cS CCManSrv) UpdateCardByID(ctx context.Context, cardID uint32, card *ccmanrpc.Card) (err error) {
	if ctx.Err() != nil {
		return cS.error(err)
	}

	in := &ccmanrpc.UpdateCardByIDReq{
		CardID:     cardID,
		UpdateCard: card,
	}

	_, err = cS.client().UpdateCardByID(ctx, in)
	if err != nil {
		return cS.error(err)
	}

	return nil
}

func (cS CCManSrv) DeleteCardByID(ctx context.Context, cardID uint32) (err error) {
	if ctx.Err() != nil {
		return cS.error(cS.error(err))
	}

	in := &ccmanrpc.DeleteCardByIDReq{
		CardID: cardID,
	}

	_, err = cS.client().DeleteCardByID(ctx, in)
	if err != nil {
		return cS.error(err)
	}

	return nil
}

func (cS CCManSrv) CreateColumn(ctx context.Context, column *ccmanrpc.Column) (err error) {
	if ctx.Err() != nil {
		return
	}

	in := &ccmanrpc.CreateColumnReq{
		CreateColumn: column,
	}

	_, err = cS.client().CreateColumn(ctx, in)
	if err != nil {
		return cS.error(err)
	}

	return nil
}

func (cS CCManSrv) GetColumnByID(ctx context.Context, colID uint32) (col *ccmanrpc.Column, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetColumnByIDReq{
		ColumnID: colID,
	}

	resp, err := cS.client().GetColumnByID(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Column, nil
}

func (cS CCManSrv) GetColumnsByTitle(ctx context.Context, title string) (cols []*ccmanrpc.Column, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetColumnsByTitleReq{
		Title: title,
	}

	resp, err := cS.client().GetColumnsByTitle(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Columns, nil
}

func (cS CCManSrv) GetColumnsByProjectID(ctx context.Context, projectID uint32) (cols []*ccmanrpc.Column, err error) {
	if ctx.Err() != nil {
		return nil, cS.error(ctx.Err())
	}

	in := &ccmanrpc.GetColumnsByProjectIDReq{
		ProjectID: projectID,
	}

	resp, err := cS.client().GetColumnsByProjectID(ctx, in)
	if err != nil {
		return nil, cS.error(err)
	}

	return resp.Columns, nil
}

func (cS CCManSrv) UpdateColumnByID(ctx context.Context, colID uint32, col *ccmanrpc.Column) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.UpdateColumnByIDReq{
		ColumnID: colID,
		Column:   col,
	}

	_, err = cS.client().UpdateColumnByID(ctx, in)
	if err != nil {
		return cS.error(err)
	}

	return nil
}

func (cS CCManSrv) DeleteColumnByID(ctx context.Context, colID uint32) (err error) {
	if ctx.Err() != nil {
		return cS.error(ctx.Err())
	}

	in := &ccmanrpc.DeleteColumnByIDReq{
		ColumnID: colID,
	}

	_, err = cS.client().DeleteColumnByID(ctx, in)
	if err != nil {
		return cS.error(err)
	}

	return nil
}
