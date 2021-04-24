package rpcserver

import (
	"context"
	"errors"

	"github.com/thanhpp/prom/pkg/logger"

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
		logger.Get().Errorf("Create card error: %v", EmptyReqError)
		return &ccmanrpc.CreateCardResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().CreateCard(ctx, req.CreateCard); err != nil {
		logger.Get().Errorf("Create card error: %v", err)
		return &ccmanrpc.CreateCardResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Create card OK")
	return &ccmanrpc.CreateCardResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) GetCardByID(ctx context.Context, req *ccmanrpc.GetCardByIDReq) (resp *ccmanrpc.GetCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get card ID error: %v", EmptyReqError)
		return &ccmanrpc.GetCardByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	card, err := repository.GetDAO().GetCardByID(ctx, req.CardID)
	if err != nil {
		logger.Get().Errorf("Get card ID error: %v", err)
		return &ccmanrpc.GetCardByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get card ID OK")
	return &ccmanrpc.GetCardByIDResp{Code: core.RPCSuccessCode, ResCard: card}, nil
}

func (c *ccManSv) GetCardsByDueDate(ctx context.Context, req *ccmanrpc.GetCardsByDueDateReq) (resp *ccmanrpc.GetCardsByDueDateResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards due date error: %v", EmptyReqError)
		return &ccmanrpc.GetCardsByDueDateResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	var duedate = timerpc.ToTime(req.DueDate)

	cards, err := repository.GetDAO().GetCardsByDueDate(ctx, duedate)
	if err != nil {
		logger.Get().Errorf("Get cards due date error %v", err)
		return &ccmanrpc.GetCardsByDueDateResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get cards due date OK")
	return &ccmanrpc.GetCardsByDueDateResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) GetCardsByAssignedToID(ctx context.Context, req *ccmanrpc.GetCardsByAssignedToIDReq) (resp *ccmanrpc.GetCardsByAssignedToIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards assigned ID error: %v", EmptyReqError)
		return &ccmanrpc.GetCardsByAssignedToIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	cards, err := repository.GetDAO().GetCardsByAssignedToID(ctx, req.AssignedToID)
	if err != nil {
		logger.Get().Errorf("Get cards assigned ID error: %v", err)
		return &ccmanrpc.GetCardsByAssignedToIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}
	logger.Get().Info("Get cards assigned ID OK")
	return &ccmanrpc.GetCardsByAssignedToIDResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) GetCardsByCreatorID(ctx context.Context, req *ccmanrpc.GetCardsByCreatorIDReq) (resp *ccmanrpc.GetCardsByCreatorIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards creator error: %v", EmptyReqError)
		return &ccmanrpc.GetCardsByCreatorIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	cards, err := repository.GetDAO().GetCardsByCreatorID(ctx, req.CreatorID)
	if err != nil {
		logger.Get().Errorf("Get cards creator error: %v", err)
		return &ccmanrpc.GetCardsByCreatorIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get cards creator OK")
	return &ccmanrpc.GetCardsByCreatorIDResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

// COLUMN

func (c *ccManSv) GetCardsByColumnID(ctx context.Context, req *ccmanrpc.GetCardsByColumnIDReq) (resp *ccmanrpc.GetCardsByColumnIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards column ID error: %v", EmptyReqError)
		return &ccmanrpc.GetCardsByColumnIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	cards, err := repository.GetDAO().GetCardsByColumnID(ctx, req.ColumnID)
	if err != nil {
		logger.Get().Errorf("Get cards column ID error: %v", err)
		return &ccmanrpc.GetCardsByColumnIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get cards column ID OK")
	return &ccmanrpc.GetCardsByColumnIDResp{Code: core.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) UpdateCardByID(ctx context.Context, req *ccmanrpc.UpdateCardByIDReq) (resp *ccmanrpc.UpdateCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Update card error : %v", EmptyReqError)
		return &ccmanrpc.UpdateCardByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().UpdateCardByID(ctx, req.CardID, req.UpdateCard); err != nil {
		logger.Get().Errorf("Update card error : %v", err)
		return &ccmanrpc.UpdateCardByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}
	logger.Get().Info("Update card OK")
	return &ccmanrpc.UpdateCardByIDResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteCardByID(ctx context.Context, req *ccmanrpc.DeleteCardByIDReq) (resp *ccmanrpc.DeleteCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Delete card error : %v", EmptyReqError)
		return &ccmanrpc.DeleteCardByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().DeleteCardByID(ctx, req.CardID); err != nil {
		logger.Get().Errorf("Delete card error : %v", err)
		return &ccmanrpc.DeleteCardByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}
	logger.Get().Info("Delete card OK")
	return &ccmanrpc.DeleteCardByIDResp{Code: core.RPCSuccessCode}, nil
}

// COLUMN

func (c *ccManSv) CreateColumn(ctx context.Context, req *ccmanrpc.CreateColumnReq) (resp *ccmanrpc.CreateColumnResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Create column error : %v", EmptyReqError)
		return &ccmanrpc.CreateColumnResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().CreateColumn(ctx, req.CreateColumn); err != nil {
		logger.Get().Errorf("Create column error : %v", err)
		return &ccmanrpc.CreateColumnResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Create column OK")
	return &ccmanrpc.CreateColumnResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) GetColumnByID(ctx context.Context, req *ccmanrpc.GetColumnByIDReq) (resp *ccmanrpc.GetColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get column by ID error : %v", EmptyReqError)
		return &ccmanrpc.GetColumnByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	column, err := repository.GetDAO().GetColumnByID(ctx, req.ColumnID)
	if err != nil {
		logger.Get().Errorf("Get column by ID error : %v", err)
		return &ccmanrpc.GetColumnByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get column by ID OK")
	return &ccmanrpc.GetColumnByIDResp{Code: core.RPCSuccessCode, Column: column}, nil
}

func (c *ccManSv) GetColumnsByTitle(ctx context.Context, req *ccmanrpc.GetColumnsByTitleReq) (resp *ccmanrpc.GetColumnsByTitleResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get column by title error : %v", EmptyReqError)
		return &ccmanrpc.GetColumnsByTitleResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	columns, err := repository.GetDAO().GetColumnsByTitle(ctx, req.Title)
	if err != nil {
		logger.Get().Errorf("Get column by title error : %v", err)
		return &ccmanrpc.GetColumnsByTitleResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get column by title OK")
	return &ccmanrpc.GetColumnsByTitleResp{Code: core.RPCSuccessCode, Columns: columns}, nil
}

func (c *ccManSv) GetColumnsByProjectID(ctx context.Context, req *ccmanrpc.GetColumnsByProjectIDReq) (resp *ccmanrpc.GetColumnsByProjectIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get column by project ID error : %v", EmptyReqError)
		return &ccmanrpc.GetColumnsByProjectIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	columns, err := repository.GetDAO().GetColumnsByProjectID(ctx, req.ProjectID)
	if err != nil {
		logger.Get().Errorf("Get column by project ID error : %v", err)
		return &ccmanrpc.GetColumnsByProjectIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Get column by project ID OK")
	return &ccmanrpc.GetColumnsByProjectIDResp{Code: core.RPCSuccessCode, Columns: columns}, nil
}

func (c *ccManSv) UpdateColumnByID(ctx context.Context, req *ccmanrpc.UpdateColumnByIDReq) (resp *ccmanrpc.UpdateColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Update column error : %v", EmptyReqError)
		return &ccmanrpc.UpdateColumnByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().UpdateColumnByID(ctx, req.ColumnID, req.Column); err != nil {
		logger.Get().Errorf("Update column error : %v", err)
		return &ccmanrpc.UpdateColumnByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Update column OK")
	return &ccmanrpc.UpdateColumnByIDResp{Code: core.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteColumnByID(ctx context.Context, req *ccmanrpc.DeleteColumnByIDReq) (resp *ccmanrpc.DeleteColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Delete column by ID error : %v", EmptyReqError)
		return &ccmanrpc.DeleteColumnByIDResp{Code: core.RPCEmptyRequestCode, Message: EmptyReqError.Error()}, EmptyReqError
	}

	if err := repository.GetDAO().DeleteColumnByID(ctx, req.ColumnID); err != nil {
		logger.Get().Errorf("Delete column by ID error : %v", err)
		return &ccmanrpc.DeleteColumnByIDResp{Code: core.DBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("Delete column by ID OK")
	return &ccmanrpc.DeleteColumnByIDResp{Code: core.RPCSuccessCode}, nil
}
