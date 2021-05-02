package rpcserver

import (
	"context"

	"github.com/thanhpp/prom/cmd/usrman/repository"
	"github.com/thanhpp/prom/pkg/errconst"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type usrManSrv struct{}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- USER ----------------------------------------------------------

func (u *usrManSrv) GetUserByID(ctx context.Context, req *usrmanrpc.GetUserByIDReq) (resp *usrmanrpc.GetUserByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetUserByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetUserByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	user, err := repository.GetDAO().GetUserByID(ctx, req.UserID)
	if err != nil {
		logger.Get().Errorf("GetUserByID error: %v", err)
		return &usrmanrpc.GetUserByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetUserByID OK")
	return &usrmanrpc.GetUserByIDResp{Code: errconst.RPCSuccessCode, User: user}, nil
}

func (u *usrManSrv) CreateUser(ctx context.Context, req *usrmanrpc.CreateUserReq) (resp *usrmanrpc.CreateUserResp, err error) {
	if req == nil {
		logger.Get().Errorf("CreateUser error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.CreateUserResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().CreateUser(ctx, req.User); err != nil {
		logger.Get().Errorf("CreateUser error: %v", err)
		return &usrmanrpc.CreateUserResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("CreateUser OK")
	return &usrmanrpc.CreateUserResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) GetUserByUsernamePass(ctx context.Context, req *usrmanrpc.GetUserByUsernamePassReq) (resp *usrmanrpc.GetUserByUsernamePassResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetUserByUsernamePass error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetUserByUsernamePassResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	user, err := repository.GetDAO().GetUserByUsernamePass(ctx, req.Username, req.Pass)
	if err != nil {
		logger.Get().Errorf("GetUserByUsernamePass error: %v", err)
		return &usrmanrpc.GetUserByUsernamePassResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetUserByUsernamePass OK")
	return &usrmanrpc.GetUserByUsernamePassResp{Code: errconst.RPCSuccessCode, User: user}, nil
}

func (u *usrManSrv) GetUserByTeamID(ctx context.Context, req *usrmanrpc.GetUserByTeamIDReq) (resp *usrmanrpc.GetUserByTeamIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetUserByTeamID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetUserByTeamIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	users, err := repository.GetDAO().GetUserByTeamID(ctx, req.TeamID)
	if err != nil {
		logger.Get().Errorf("GetUserByTeamID error: %v", err)
		return &usrmanrpc.GetUserByTeamIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetUserByTeamID OK")
	return &usrmanrpc.GetUserByTeamIDResp{Code: errconst.RPCSuccessCode, Users: users}, nil
}

func (u *usrManSrv) UpdateUserByID(ctx context.Context, req *usrmanrpc.UpdateUserByIDReq) (resp *usrmanrpc.UpdateUserByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("UpdateUserByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.UpdateUserByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().UpdateUserByID(ctx, req.UserID, req.User); err != nil {
		logger.Get().Errorf("UpdateUserByID error: %v", err)
		return &usrmanrpc.UpdateUserByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("UpdateUserByID OK")
	return &usrmanrpc.UpdateUserByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) DeleteUserByID(ctx context.Context, req *usrmanrpc.DeleteUserByIDReq) (resp *usrmanrpc.DeleteUserByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("DeleteUserByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.DeleteUserByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().DeleteUserByID(ctx, req.UserID); err != nil {
		logger.Get().Errorf("DeleteUserByID error: %v", err)
		return &usrmanrpc.DeleteUserByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("DeleteUserByID OK")
	return &usrmanrpc.DeleteUserByIDResp{Code: errconst.RPCSuccessCode}, nil
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- TEAM ----------------------------------------------------------

func (u *usrManSrv) CreateTeam(ctx context.Context, req *usrmanrpc.CreateTeamReq) (resp *usrmanrpc.CreateTeamResp, err error) {
	if req == nil {
		logger.Get().Errorf("CreateTeam error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.CreateTeamResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().CreateTeam(ctx, req.Team); err != nil {
		logger.Get().Errorf("CreateTeam error: %v", err)
		return &usrmanrpc.CreateTeamResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("CreateTeam OK")
	return &usrmanrpc.CreateTeamResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) GetTeamByID(ctx context.Context, req *usrmanrpc.GetTeamByIDReq) (resp *usrmanrpc.GetTeamByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetTeamByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetTeamByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	team, err := repository.GetDAO().GetTeamByID(ctx, req.TeamID)
	if err != nil {
		logger.Get().Errorf("GetTeamByID error: %v", err)
		return &usrmanrpc.GetTeamByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetTeamByID OK")
	return &usrmanrpc.GetTeamByIDResp{Code: errconst.RPCSuccessCode, Team: team}, nil
}

func (u *usrManSrv) GetTeamsByUserID(ctx context.Context, req *usrmanrpc.GetTeamsByUserIDReq) (resp *usrmanrpc.GetTeamsByUserIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetTeamsByUserID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetTeamsByUserIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	teams, err := repository.GetDAO().GetTeamsByUserID(ctx, req.UserID)
	if err != nil {
		logger.Get().Errorf("GetTeamsByUserID error: %v", err)
		return &usrmanrpc.GetTeamsByUserIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetTeamsByUserID OK")
	return &usrmanrpc.GetTeamsByUserIDResp{Code: errconst.RPCSuccessCode, Teams: teams}, nil
}

func (u *usrManSrv) GetTeamsByCreatorID(ctx context.Context, req *usrmanrpc.GetTeamsByCreatorIDReq) (resp *usrmanrpc.GetTeamsByCreatorIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetTeamsByCreatorID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetTeamsByCreatorIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	teams, err := repository.GetDAO().GetTeamsByCreatorID(ctx, req.CreatorID)
	if err != nil {
		logger.Get().Errorf("GetTeamsByCreatorID error: %v", err)
		return &usrmanrpc.GetTeamsByCreatorIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetTeamsByCreatorID OK")
	return &usrmanrpc.GetTeamsByCreatorIDResp{Code: errconst.RPCSuccessCode, Teams: teams}, nil
}

func (u *usrManSrv) GetTeamByName(ctx context.Context, req *usrmanrpc.GetTeamByNameReq) (resp *usrmanrpc.GetTeamByNameResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetTeamByName error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetTeamByNameResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	teams, err := repository.GetDAO().GetTeamByName(ctx, req.TeamName)
	if err != nil {
		logger.Get().Errorf("GetTeamByName error: %v", err)
		return &usrmanrpc.GetTeamByNameResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetTeamByName OK")
	return &usrmanrpc.GetTeamByNameResp{Code: errconst.RPCSuccessCode, Teams: teams}, nil
}

func (u *usrManSrv) UpdateTeamByID(ctx context.Context, req *usrmanrpc.UpdateTeamByIDReq) (resp *usrmanrpc.UpdateTeamByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("UpdateTeamByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.UpdateTeamByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().UpdateTeamByID(ctx, req.TeamID, req.Team); err != nil {
		logger.Get().Errorf("UpdateTeamByID error: %v", err)
		return &usrmanrpc.UpdateTeamByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("UpdateTeamByID OK")
	return &usrmanrpc.UpdateTeamByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) AddMemberByID(ctx context.Context, req *usrmanrpc.AddMemberByIDReq) (resp *usrmanrpc.AddMemberByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("AddMemberByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.AddMemberByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().AddMemberByID(ctx, req.TeamID, req.TeamID); err != nil {
		logger.Get().Errorf("AddMemberByID error: %v", err)
		return &usrmanrpc.AddMemberByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("AddMemberByID OK")
	return &usrmanrpc.AddMemberByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) RemoveMemberByID(ctx context.Context, req *usrmanrpc.RemoveMemberByIDReq) (resp *usrmanrpc.RemoveMemberByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("RemoveMemberByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.RemoveMemberByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().RemoveMemberByID(ctx, req.TeamID, req.UserID); err != nil {
		logger.Get().Errorf("RemoveMemberByID error: %v", err)
		return &usrmanrpc.RemoveMemberByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("RemoveMemberByID OK")
	return &usrmanrpc.RemoveMemberByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) DeleteTeamByID(ctx context.Context, req *usrmanrpc.DeleteTeamByIDReq) (resp *usrmanrpc.DeleteTeamByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("DeleteTeamByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.DeleteTeamByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().DeleteTeamByID(ctx, req.TeamID); err != nil {
		logger.Get().Errorf("DeleteTeamByID error: %v", err)
		return &usrmanrpc.DeleteTeamByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("DeleteTeamByID OK")
	return &usrmanrpc.DeleteTeamByIDResp{Code: errconst.RPCSuccessCode}, nil
}

// ---------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- PROJECT ----------------------------------------------------------

func (u *usrManSrv) CreateProject(ctx context.Context, req *usrmanrpc.CreateProjectReq) (resp *usrmanrpc.CreateProjectResp, err error) {
	if req == nil {
		logger.Get().Errorf("CreateProject error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.CreateProjectResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().CreateProject(ctx, req.Project); err != nil {
		logger.Get().Errorf("CreateProject error: %v", err)
		return &usrmanrpc.CreateProjectResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("CreateProject OK")
	return &usrmanrpc.CreateProjectResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) GetProjectByID(ctx context.Context, req *usrmanrpc.GetProjectByIDReq) (resp *usrmanrpc.GetProjectByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetProjectByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetProjectByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	project, err := repository.GetDAO().GetProjectByID(ctx, req.ProjectID)
	if err != nil {
		logger.Get().Errorf("GetProjectByID error: %v", err)
		return &usrmanrpc.GetProjectByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetProjectByID OK")
	return &usrmanrpc.GetProjectByIDResp{Code: errconst.RPCSuccessCode, Project: project}, nil
}

func (u *usrManSrv) GetProjtectsByTeamID(ctx context.Context, req *usrmanrpc.GetProjtectsByTeamIDReq) (resp *usrmanrpc.GetProjtectsByTeamIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("GetProjtectsByTeamID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.GetProjtectsByTeamIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	projects, err := repository.GetDAO().GetProjtectsByTeamID(ctx, req.TeamID)
	if err != nil {
		logger.Get().Errorf("GetProjtectsByTeamID error: %v", err)
		return &usrmanrpc.GetProjtectsByTeamIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("GetProjtectsByTeamID OK")
	return &usrmanrpc.GetProjtectsByTeamIDResp{Code: errconst.RPCSuccessCode, Projects: projects}, nil
}

func (u *usrManSrv) UpdateProjectByID(ctx context.Context, req *usrmanrpc.UpdateProjectByIDReq) (resp *usrmanrpc.UpdateProjectByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("UpdateProjectByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.UpdateProjectByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().UpdateProjectByID(ctx, req.ProjectID, req.Project); err != nil {
		logger.Get().Errorf("UpdateProjectByID error: %v", err)
		return &usrmanrpc.UpdateProjectByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("UpdateProjectByID OK")
	return &usrmanrpc.UpdateProjectByIDResp{Code: errconst.RPCSuccessCode}, nil
}

func (u *usrManSrv) DeleteProjectByID(ctx context.Context, req *usrmanrpc.DeleteProjectByIDReq) (resp *usrmanrpc.DeleteProjectByIDResp, err error) {
	if req == nil {
		logger.Get().Errorf("DeleteProjectByID error: %v", errconst.RPCEmptyRequestErr)
		return &usrmanrpc.DeleteProjectByIDResp{Code: errconst.RPCEmptyReqCode}, errconst.RPCEmptyRequestErr
	}

	if err = repository.GetDAO().DeleteProjectByID(ctx, req.ProjectID); err != nil {
		logger.Get().Errorf("DeleteProjectByID error: %v", err)
		return &usrmanrpc.DeleteProjectByIDResp{Code: errconst.RPCDBErrorCode, Message: err.Error()}, err
	}

	logger.Get().Info("DeleteProjectByID OK")
	return &usrmanrpc.DeleteProjectByIDResp{Code: errconst.RPCSuccessCode}, nil
}