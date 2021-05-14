package repository

import (
	"context"

	"github.com/thanhpp/prom/cmd/usrman/repository/gormdb"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type iDAO interface {
	InitDBConnection(dsn string, logLevel string) (err error)
	AutoMigrate(ctx context.Context, models ...interface{}) (err error)

	CreateUser(ctx context.Context, usr *usrmanrpc.User) (err error)
	GetUserByID(ctx context.Context, usrID uint32) (usr *usrmanrpc.User, err error)
	GetUserByUsernamePass(ctx context.Context, usrname string, hashpwd string) (usr *usrmanrpc.User, err error)
	GetUserByTeamID(ctx context.Context, teamID uint32) (usrs []*usrmanrpc.User, err error)
	UpdateUserByID(ctx context.Context, usrID uint32, usr *usrmanrpc.User) (err error)
	DeleteUserByID(ctx context.Context, usrID uint32) (err error)

	CreateTeam(ctx context.Context, team *usrmanrpc.Team) (err error)
	GetTeamByID(ctx context.Context, teamID uint32) (team *usrmanrpc.Team, err error)
	GetTeamsByUserID(ctx context.Context, userID uint32) (teams []*usrmanrpc.Team, err error)
	GetTeamsByCreatorID(ctx context.Context, creatorID uint32) (teams []*usrmanrpc.Team, err error)
	GetTeamByName(ctx context.Context, name string) (teams []*usrmanrpc.Team, err error)
	UpdateTeamByID(ctx context.Context, teamID uint32, team *usrmanrpc.Team) (err error)
	AddMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error)
	RemoveMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error)
	DeleteTeamByID(ctx context.Context, teamID uint32) (err error)

	CreateProject(ctx context.Context, project *usrmanrpc.Project) (err error)
	GetProjectByID(ctx context.Context, projectID uint32) (project *usrmanrpc.Project, err error)
	GetProjtectsByTeamID(ctx context.Context, teamID uint32) (projects []*usrmanrpc.Project, err error)
	UpdateProjectByID(ctx context.Context, projectID uint32, project *usrmanrpc.Project) (err error)
	AddColumnsToProject(ctx context.Context, projectID uint32, columnID uint32) (err error)
	DeleteProjectByID(ctx context.Context, projectID uint32) (err error)
}

func GetDAO() iDAO {
	return gormdb.GetGormDB()
}
