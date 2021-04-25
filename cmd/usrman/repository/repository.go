package repository

import (
	"context"

	"github.com/thanhpp/prom/cmd/usrman/repository/entity"
	"github.com/thanhpp/prom/cmd/usrman/repository/gormdb"
)

type iDAO interface {
	InitDBConnection(dsn string, logLevel string) (err error)
	AutoMigrate(ctx context.Context, models ...interface{}) (err error)

	CreateUser(ctx context.Context, usr *entity.User) (err error)
	GetUserByID(ctx context.Context, usrID uint32) (usr *entity.User, err error)
	GetUserByUsernamePass(ctx context.Context, usrname string, hashpwd string) (usr *entity.User, err error)
	GetUserByTeamID(ctx context.Context, teamID uint32) (usrs []*entity.User, err error)
	UpdateUserByID(ctx context.Context, usrID uint32, usr *entity.User) (err error)
	DeleteUserByID(ctx context.Context, usrID uint32) (err error)

	CreateTeam(ctx context.Context, team *entity.Team) (err error)
	GetTeamByID(ctx context.Context, teamID *entity.Team) (team *entity.Team, err error)
	GetTeamsByCreatorID(ctx context.Context, creatorID uint32) (teams []*entity.Team, err error)
	GetTeamByName(ctx context.Context, name string) (teams []*entity.Team, err error)
	UpdateTeamByID(ctx context.Context, teamID uint32, team *entity.Team) (err error)
	AddMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error)
	RemoveMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error)
	DeleteTeamByID(ctx context.Context, teamID uint32) (err error)

	CreateProject(ctx context.Context, project *entity.Project) (err error)
	GetProjectByID(ctx context.Context, projectID uint32) (project *entity.Project, err error)
	GetProjtectsByTeamID(ctx context.Context, teamID uint32) (projects []*entity.Project, err error)
	UpdateProjectByID(ctx context.Context, projectID uint32, project *entity.Project) (err error)
	DeleteProjectByID(ctx context.Context, projectID uint32) (err error)
}

func GetDAO() iDAO {
	return gormdb.GetGormDB()
}
