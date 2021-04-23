package rpcserver

import (
	"context"
	"errors"

	"github.com/thanhpp/prom/pkg/timerpc"

	"github.com/thanhpp/prom/cmd/ccman/core"
	"github.com/thanhpp/prom/cmd/ccman/repository"

	"github.com/thanhpp/prom/pkg/ccmanrpc"
)

// ------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- RPC SERVER ----------------------------------------------------------

type ccManSv struct{}

type iCCManSv interface {
	CreateCard(ctx context.Context, req *ccmanrpc.CreateCardReq) (resp *ccmanrpc.CreateCardResp, err error)
	GetCardByID(ctx context.Context, req *ccmanrpc.GetCardByIDReq) (resp *ccmanrpc.GetCardByIDResp, err error)
	GetCardsByDueDate(ctx context.Context, req *ccmanrpc.GetCardsByDueDateReq) (resp *ccmanrpc.GetCardsByDueDateResp, err error)
	GetCardsByAssignedToID(ctx context.Context, req *ccmanrpc.GetCardsByAssignedToIDReq) (resp *ccmanrpc.GetCardsByAssignedToIDResp, err error)
	GetCardsByCreatorID(ctx context.Context, req *ccmanrpc.GetCardsByCreatorIDReq) (resp *ccmanrpc.GetCardsByCreatorIDResp, err error)
	GetCardsByColumnID(ctx context.Context, req *ccmanrpc.GetCardsByColumnIDReq) (resp *ccmanrpc.GetCardsByColumnIDResp, err error)
	UpdateCardByID(ctx context.Context, req *ccmanrpc.UpdateCardByIDReq) (resp *ccmanrpc.UpdateCardByIDResp, err error)
	DeleteCardByID(ctx context.Context, req *ccmanrpc.DeleteCardByIDReq) (resp *ccmanrpc.DeleteCardByIDResp, err error)
	CreateColumn(ctx context.Context, req *ccmanrpc.CreateColumnReq) (resp *ccmanrpc.CreateColumnResp, err error)
	GetColumnByID(ctx context.Context, req *ccmanrpc.GetColumnByIDReq) (resp *ccmanrpc.GetColumnByIDResp, err error)
	GetColumnsByTitle(ctx context.Context, req *ccmanrpc.GetColumnsByTitleReq) (resp *ccmanrpc.GetColumnsByTitleResp, err error)
	GetColumnsByProjectID(ctx context.Context, req *ccmanrpc.GetColumnsByProjectIDReq) (resp *ccmanrpc.GetColumnsByProjectIDResp, err error)
	UpdateColumnByID(ctx context.Context, req *ccmanrpc.UpdateColumnByIDReq) (resp *ccmanrpc.UpdateColumnByIDResp, err error)
	DeleteColumnByID(ctx context.Context, req *ccmanrpc.DeleteColumnByIDReq) (resp *ccmanrpc.DeleteColumnByIDResp, err error)
}

var _ iCCManSv = (*ccManSv)(nil) //compile check

var EmptyReqError = errors.New("Empty RPC request")

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- IMPLEMENT ----------------------------------------------------------

// CARD

func (c *ccManSv) CreateCard(ctx context.Context, req *ccmanrpc.CreateCardReq) (resp *ccmanrpc.CreateCardResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.CreateCardResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().CreateCard(ctx, req.CreateCard); err != nil {
		return &ccmanrpc.CreateCardResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.CreateCardResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) GetCardByID(ctx context.Context, req *ccmanrpc.GetCardByIDReq) (resp *ccmanrpc.GetCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetCardByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	card, err := repository.GetDAO().GetCardByID(ctx, req.CardID)
	if err != nil {
		return &ccmanrpc.GetCardByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetCardByIDResp{Code: core.RPCSuccessCode, ResCard: card}, nil
}

func (c *ccManSv) GetCardsByDueDate(ctx context.Context, req *ccmanrpc.GetCardsByDueDateReq) (resp *ccmanrpc.GetCardsByDueDateResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetCardsByDueDateResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	var duedate = timerpc.ToTime(req.DueDate)

	cards, err := repository.GetDAO().GetCardsByDueDate(ctx, duedate)
	if err != nil {
		return &ccmanrpc.GetCardsByDueDateResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetCardsByDueDateResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) GetCardsByAssignedToID(ctx context.Context, req *ccmanrpc.GetCardsByAssignedToIDReq) (resp *ccmanrpc.GetCardsByAssignedToIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetCardsByAssignedToIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	cards, err := repository.GetDAO().GetCardsByAssignedToID(ctx, req.AssignedToID)
	if err != nil {
		return &ccmanrpc.GetCardsByAssignedToIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetCardsByAssignedToIDResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) GetCardsByCreatorID(ctx context.Context, req *ccmanrpc.GetCardsByCreatorIDReq) (resp *ccmanrpc.GetCardsByCreatorIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetCardsByCreatorIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	cards, err := repository.GetDAO().GetCardsByCreatorID(ctx, req.CreatorID)
	if err != nil {
		return &ccmanrpc.GetCardsByCreatorIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetCardsByCreatorIDResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

// COLUMN

func (c *ccManSv) GetCardsByColumnID(ctx context.Context, req *ccmanrpc.GetCardsByColumnIDReq) (resp *ccmanrpc.GetCardsByColumnIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetCardsByColumnIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	cards, err := repository.GetDAO().GetCardsByColumnID(ctx, req.ColumnID)
	if err != nil {
		return &ccmanrpc.GetCardsByColumnIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetCardsByColumnIDResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) UpdateCardByID(ctx context.Context, req *ccmanrpc.UpdateCardByIDReq) (resp *ccmanrpc.UpdateCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.UpdateCardByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().UpdateCardByID(ctx, req.CardID, req.UpdateCard); err != nil {
		return &ccmanrpc.UpdateCardByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.UpdateCardByIDResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteCardByID(ctx context.Context, req *ccmanrpc.DeleteCardByIDReq) (resp *ccmanrpc.DeleteCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.DeleteCardByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().DeleteCardByID(ctx, req.CardID); err != nil {
		return &ccmanrpc.DeleteCardByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.DeleteCardByIDResp{Code: core.RPCSuccessCode}, nil
}

// COLUMN

func (c *ccManSv) CreateColumn(ctx context.Context, req *ccmanrpc.CreateColumnReq) (resp *ccmanrpc.CreateColumnResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.CreateColumnResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().CreateColumn(ctx, req.CreateColumn); err != nil {
		return &ccmanrpc.CreateColumnResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.CreateColumnResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) GetColumnByID(ctx context.Context, req *ccmanrpc.GetColumnByIDReq) (resp *ccmanrpc.GetColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetColumnByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	column, err := repository.GetDAO().GetColumnByID(ctx, req.ColumnID)
	if err != nil {
		return &ccmanrpc.GetColumnByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetColumnByIDResp{Code: core.RPCSuccessCode, Column: column}, nil
}

func (c *ccManSv) GetColumnsByTitle(ctx context.Context, req *ccmanrpc.GetColumnsByTitleReq) (resp *ccmanrpc.GetColumnsByTitleResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetColumnsByTitleResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	columns, err := repository.GetDAO().GetColumnsByTitle(ctx, req.Title)
	if err != nil {
		return &ccmanrpc.GetColumnsByTitleResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetColumnsByTitleResp{Code: core.RPCSuccessCode, Columns: columns}, nil
}

func (c *ccManSv) GetColumnsByProjectID(ctx context.Context, req *ccmanrpc.GetColumnsByProjectIDReq) (resp *ccmanrpc.GetColumnsByProjectIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.GetColumnsByProjectIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	columns, err := repository.GetDAO().GetColumnsByProjectID(ctx, req.ProjectID)
	if err != nil {
		return &ccmanrpc.GetColumnsByProjectIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.GetColumnsByProjectIDResp{Code: core.RPCSuccessCode, Columns: columns}, nil
}

func (c *ccManSv) UpdateColumnByID(ctx context.Context, req *ccmanrpc.UpdateColumnByIDReq) (resp *ccmanrpc.UpdateColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.UpdateColumnByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().UpdateColumnByID(ctx, req.ColumnID, req.Column); err != nil {
		return &ccmanrpc.UpdateColumnByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.UpdateColumnByIDResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteColumnByID(ctx context.Context, req *ccmanrpc.DeleteColumnByIDReq) (resp *ccmanrpc.DeleteColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		return &ccmanrpc.DeleteColumnByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().DeleteColumnByID(ctx, req.ColumnID); err != nil {
		return &ccmanrpc.DeleteColumnByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	return &ccmanrpc.DeleteColumnByIDResp{Code: core.RPCSuccessCode}, nil
}
