package service

import (
	"context"

	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type UsrManSrv struct {
}

type iUsrManSrv interface {
	// user
	Login(ctx context.Context, username string, pass string) (jwt string, err error)
	Logout(ctx context.Context, userID uint32) (err error)
	NewUser(ctx context.Context, username string, pass string) (err error)
	UpdateUsername(ctx context.Context, userID uint32, username string) (err error)
	UpdatePassword(ctx context.Context, userID uint32, password string) (err error)

	// team
	GetTeamsByUserID(ctx context.Context, userID uint32) (teams []*usrmanrpc.Team, err error)
	GetTeamMembersByID(ctx context.Context, teamID uint32) (users []*usrmanrpc.User, err error)
	AddMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error)
	RemoveMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error)

	// project
	GetProjectsByTeamID(ctx context.Context, teamID uint32) (projects []*usrmanrpc.Project, err error)
	NewProject(ctx context.Context, project *usrmanrpc.Project) (err error)
	UpdateProject(ctx context.Context, project *usrmanrpc.Project) (err error)
	DeleteProjectByID(ctx context.Context, projectID uint32) (err error)
}
