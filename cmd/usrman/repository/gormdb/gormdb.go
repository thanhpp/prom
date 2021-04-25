package gormdb

import (
	"context"
	"log"

	"github.com/thanhpp/prom/cmd/usrman/repository/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// --------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- GORMDB ----------------------------------------------------------

type implGorm struct{}

var (
	gDB       = &gorm.DB{}
	gormObj   = new(implGorm)
	usrModel  = new(entity.User)
	teamModel = new(entity.Team)
	prjModel  = new(entity.Project)
)

// GetGormDB ...
func GetGormDB() *implGorm {
	return gormObj
}

// ---------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- IMPLEMENT DAO ----------------------------------------------------------

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

func (g *implGorm) CreateUser(ctx context.Context, usr *entity.User) (err error) {
	return
}

func (g *implGorm) GetUserByID(ctx context.Context, usrID uint32) (usr *entity.User, err error) {
	return
}

func (g *implGorm) GetUserByUsernamePass(ctx context.Context, usrname string, hashpwd string) (usr *entity.User, err error) {
	return
}

func (g *implGorm) GetUserByTeamID(ctx context.Context, teamID uint32) (usrs []*entity.User, err error) {
	return
}

func (g *implGorm) UpdateUserByID(ctx context.Context, usrID uint32, usr *entity.User) (err error) {
	return
}

func (g *implGorm) DeleteUserByID(ctx context.Context, usrID uint32) (err error) {
	return
}

func (g *implGorm) CreateTeam(ctx context.Context, team *entity.Team) (err error) {
	return
}

func (g *implGorm) GetTeamByID(ctx context.Context, teamID *entity.Team) (team *entity.Team, err error) {
	return
}

func (g *implGorm) GetTeamsByCreatorID(ctx context.Context, creatorID uint32) (teams []*entity.Team, err error) {
	return
}

func (g *implGorm) GetTeamByName(ctx context.Context, name string) (teams []*entity.Team, err error) {
	return
}

func (g *implGorm) UpdateTeamByID(ctx context.Context, teamID uint32, team *entity.Team) (err error) {
	return
}

func (g *implGorm) AddMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error) {
	return
}

func (g *implGorm) RemoveMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error) {
	return
}

func (g *implGorm) DeleteTeamByID(ctx context.Context, teamID uint32) (err error) {
	return
}

func (g *implGorm) CreateProject(ctx context.Context, project *entity.Project) (err error) {
	return
}

func (g *implGorm) GetProjectByID(ctx context.Context, projectID uint32) (project *entity.Project, err error) {
	return
}

func (g *implGorm) GetProjtectsByTeamID(ctx context.Context, teamID uint32) (projects []*entity.Project, err error) {
	return
}

func (g *implGorm) UpdateProjectByID(ctx context.Context, projectID uint32, project *entity.Project) (err error) {
	return
}

func (g *implGorm) DeleteProjectByID(ctx context.Context, projectID uint32) (err error) {
	return
}
