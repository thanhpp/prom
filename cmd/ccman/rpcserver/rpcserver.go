package rpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/thanhpp/prom/cmd/ccman/repository"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/errconst"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/timerpc"
)

// ------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- RPC SERVER ----------------------------------------------------------

type ccManSv struct{}

type iCCManSv interface {
	GetAllFromProjectID(ctx context.Context, req *ccmanrpc.GetAllFromProjectIDReq) (resp *ccmanrpc.GetAllFromProjectIDResp, err error)
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
	DeleteColumnByIDAndMove(ctx context.Context, req *ccmanrpc.DeleteColumnByIDAndMoveReq) (resp *ccmanrpc.DeleteColumnByIDAndMoveResp, err error)
}

var _ iCCManSv = (*ccManSv)(nil) //compile check

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- IMPLEMENT ----------------------------------------------------------

// GET ALL

func (c *ccManSv) GetAllFromProjectID(ctx context.Context, req *ccmanrpc.GetAllFromProjectIDReq) (resp *ccmanrpc.GetAllFromProjectIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("Create card error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	cols, err := repository.GetDAO().GetAllFromProjectID(ctx, req.ProjectID)
	if err != nil {
		logger.Get().Errorf("GetAllFromProjectID error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("GetAllFromProjectID OK")

	return &ccmanrpc.GetAllFromProjectIDResp{Code: errconst.RPCSuccessCode, Columns: cols}, nil
}

// CARD

func (c *ccManSv) CreateCard(ctx context.Context, req *ccmanrpc.CreateCardReq) (resp *ccmanrpc.CreateCardResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Create card error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	id, err := repository.GetDAO().CreateCard(ctx, req.CreateCard)
	if err != nil {
		logger.Get().Errorf("Create card error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Create card OK")
	return &ccmanrpc.CreateCardResp{Code: errconst.RPCSuccessCode, CreatedID: id}, nil
}

func (c *ccManSv) GetCardByID(ctx context.Context, req *ccmanrpc.GetCardByIDReq) (resp *ccmanrpc.GetCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get card ID error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	card, err := repository.GetDAO().GetCardByID(ctx, req.CardID)
	if err != nil {
		logger.Get().Errorf("Get card ID error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get card ID OK")
	return &ccmanrpc.GetCardByIDResp{Code: errconst.RPCSuccessCode, ResCard: card}, nil
}

func (c *ccManSv) GetCardsByDueDate(ctx context.Context, req *ccmanrpc.GetCardsByDueDateReq) (resp *ccmanrpc.GetCardsByDueDateResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards due date error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	var duedate = timerpc.ToTime(req.DueDate)

	cards, err := repository.GetDAO().GetCardsByDueDate(ctx, duedate)
	if err != nil {
		logger.Get().Errorf("Get cards due date error %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get cards due date OK")
	return &ccmanrpc.GetCardsByDueDateResp{Code: errconst.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) GetCardsByAssignedToID(ctx context.Context, req *ccmanrpc.GetCardsByAssignedToIDReq) (resp *ccmanrpc.GetCardsByAssignedToIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards assigned ID error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	cards, err := repository.GetDAO().GetCardsByAssignedToID(ctx, req.AssignedToID)
	if err != nil {
		logger.Get().Errorf("Get cards assigned ID error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	logger.Get().Info("Get cards assigned ID OK")
	return &ccmanrpc.GetCardsByAssignedToIDResp{Code: errconst.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) GetCardsByCreatorID(ctx context.Context, req *ccmanrpc.GetCardsByCreatorIDReq) (resp *ccmanrpc.GetCardsByCreatorIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards creator error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	cards, err := repository.GetDAO().GetCardsByCreatorID(ctx, req.CreatorID)
	if err != nil {
		logger.Get().Errorf("Get cards creator error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get cards creator OK")
	return &ccmanrpc.GetCardsByCreatorIDResp{Code: errconst.RPCSuccessCode, Cards: cards}, nil
}

// COLUMN

func (c *ccManSv) GetCardsByColumnID(ctx context.Context, req *ccmanrpc.GetCardsByColumnIDReq) (resp *ccmanrpc.GetCardsByColumnIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get cards column ID error: %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	cards, err := repository.GetDAO().GetCardsByColumnID(ctx, req.ColumnID)
	if err != nil {
		logger.Get().Errorf("Get cards column ID error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get cards column ID OK")
	return &ccmanrpc.GetCardsByColumnIDResp{Code: errconst.RPCSuccessCode, Cards: cards}, nil
}

func (c *ccManSv) UpdateCardByID(ctx context.Context, req *ccmanrpc.UpdateCardByIDReq) (resp *ccmanrpc.UpdateCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Update card error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	if err := repository.GetDAO().UpdateCardByID(ctx, req.CardID, req.UpdateCard); err != nil {
		logger.Get().Errorf("Update card error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	logger.Get().Info("Update card OK")
	return &ccmanrpc.UpdateCardByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteCardByID(ctx context.Context, req *ccmanrpc.DeleteCardByIDReq) (resp *ccmanrpc.DeleteCardByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Delete card error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	if err := repository.GetDAO().DeleteCardByID(ctx, req.CardID); err != nil {
		logger.Get().Errorf("Delete card error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	logger.Get().Info("Delete card OK")
	return &ccmanrpc.DeleteCardByIDResp{Code: errconst.RPCSuccessCode}, nil
}

// COLUMN

func (c *ccManSv) CreateColumn(ctx context.Context, req *ccmanrpc.CreateColumnReq) (resp *ccmanrpc.CreateColumnResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Create column error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	id, err := repository.GetDAO().CreateColumn(ctx, req.CreateColumn)
	if err != nil {
		logger.Get().Errorf("Create column error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Create column OK")
	return &ccmanrpc.CreateColumnResp{Code: errconst.RPCSuccessCode, CreatedID: id}, nil
}

func (c *ccManSv) GetColumnByID(ctx context.Context, req *ccmanrpc.GetColumnByIDReq) (resp *ccmanrpc.GetColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get column by ID error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	column, err := repository.GetDAO().GetColumnByID(ctx, req.ColumnID)
	if err != nil {
		logger.Get().Errorf("Get column by ID error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get column by ID OK")
	return &ccmanrpc.GetColumnByIDResp{Code: errconst.RPCSuccessCode, Column: column}, nil
}

func (c *ccManSv) GetColumnsByTitle(ctx context.Context, req *ccmanrpc.GetColumnsByTitleReq) (resp *ccmanrpc.GetColumnsByTitleResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get column by title error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	columns, err := repository.GetDAO().GetColumnsByTitle(ctx, req.Title)
	if err != nil {
		logger.Get().Errorf("Get column by title error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get column by title OK")
	return &ccmanrpc.GetColumnsByTitleResp{Code: errconst.RPCSuccessCode, Columns: columns}, nil
}

func (c *ccManSv) GetColumnsByProjectID(ctx context.Context, req *ccmanrpc.GetColumnsByProjectIDReq) (resp *ccmanrpc.GetColumnsByProjectIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Get column by project ID error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	columns, err := repository.GetDAO().GetColumnsByProjectID(ctx, req.ProjectID)
	if err != nil {
		logger.Get().Errorf("Get column by project ID error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Get column by project ID OK")
	return &ccmanrpc.GetColumnsByProjectIDResp{Code: errconst.RPCSuccessCode, Columns: columns}, nil
}

func (c *ccManSv) UpdateColumnByID(ctx context.Context, req *ccmanrpc.UpdateColumnByIDReq) (resp *ccmanrpc.UpdateColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Update column error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	if err := repository.GetDAO().UpdateColumnByID(ctx, req.ColumnID, req.Column); err != nil {
		logger.Get().Errorf("Update column error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Update column OK")
	return &ccmanrpc.UpdateColumnByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteColumnByID(ctx context.Context, req *ccmanrpc.DeleteColumnByIDReq) (resp *ccmanrpc.DeleteColumnByIDResp, err error) {
	// pre-exec check
	if req == nil {
		logger.Get().Errorf("Delete column by ID error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	if err := repository.GetDAO().DeleteColumnByID(ctx, req.ColumnID); err != nil {
		logger.Get().Errorf("Delete column by ID error : %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("Delete column by ID OK")
	return &ccmanrpc.DeleteColumnByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (c *ccManSv) DeleteColumnByIDAndMove(ctx context.Context, req *ccmanrpc.DeleteColumnByIDAndMoveReq) (resp *ccmanrpc.DeleteColumnByIDAndMoveResp, err error) {
	if req == nil {
		logger.Get().Errorf("Delete column by ID error : %v", errconst.RPCEmptyRequestErr)
		return nil, status.Error(codes.Unavailable, "Empty request")
	}

	if err = repository.GetDAO().DeleteColumnByIDAndMove(ctx, req.ColumnID, req.NewColumnID); err != nil {
		logger.Get().Errorf("Delete & move column error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Get().Info("DeleteColumnByIDAndMove OK")
	return &ccmanrpc.DeleteColumnByIDAndMoveResp{Code: errconst.RPCSuccessCode}, nil
}
