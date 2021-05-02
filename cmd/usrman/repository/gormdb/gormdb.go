package gormdb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
	var gormLogConfig = gormlog.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormlog.Info,
		Colorful:      false,
	}

	// create db log
	switch logLevel {
	case "WARN":
		gormLogConfig.LogLevel = gormlog.Warn
	case "ERROR":
		gormLogConfig.LogLevel = gormlog.Error
	case "SILENT":
		gormLogConfig.LogLevel = gormlog.Silent
	default:
		log.Println("START GORM LOG WITH DEFAULT CONFIG: INFO")
	}

	var (
		gormConfig = &gorm.Config{
			Logger: gormlog.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormLogConfig),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
	)

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
	if err = gDB.Model(usrModel).WithContext(ctx).Save(usr).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) GetUserByID(ctx context.Context, usrID uint32) (usr *entity.User, err error) {
	usr = new(entity.User)
	if err = gDB.Model(usrModel).WithContext(ctx).Where("id = ?", usrID).Take(usr).Error; err != nil {
		return nil, err
	}

	return usr, nil
}

func (g *implGorm) GetUserByUsernamePass(ctx context.Context, usrname string, hashpwd string) (usr *entity.User, err error) {
	usr = new(entity.User)
	if err = gDB.Model(usrModel).WithContext(ctx).Where("username LIKE ? AND hash_pass LIKE ?", usrname, hashpwd).Error; err != nil {
		return nil, err
	}

	return usr, nil
}

func (g *implGorm) GetUserByTeamID(ctx context.Context, teamID uint32) (usrs []*entity.User, err error) {
	if err = gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var usersID []int

		if err := tx.WithContext(ctx).Table("team_user").Select("user_id").Where("team_id = ?", teamID).Pluck("user_id", &usersID).Error; err != nil {
			return err
		}

		rows2, err2 := tx.Model(usrModel).WithContext(ctx).Where("id IN ?", usersID).Rows()
		if err2 != nil {
			return err2
		}

		usrs, err = scanUsers(tx, rows2)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return usrs, nil
}

func (g *implGorm) UpdateUserByID(ctx context.Context, usrID uint32, usr *entity.User) (err error) {
	if err = gDB.Model(usrModel).WithContext(ctx).Where("id = ?", usrID).Updates(usr).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) DeleteUserByID(ctx context.Context, usrID uint32) (err error) {
	if err = gDB.Model(usrModel).WithContext(ctx).Where("id = ?", usrID).Delete(entity.User{}).Error; err != nil {
		return err
	}

	return nil
}

// TEAM

func (g *implGorm) CreateTeam(ctx context.Context, team *entity.Team) (err error) {
	if err = gDB.Model(teamModel).WithContext(ctx).Save(team).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) GetTeamByID(ctx context.Context, teamID uint32) (team *entity.Team, err error) {
	team = new(entity.Team)
	if err = gDB.Model(teamModel).WithContext(ctx).Where("id = ?", teamID).Take(team).Error; err != nil {
		return nil, err
	}

	return team, nil
}

func (g *implGorm) GetTeamsByUserID(ctx context.Context, userID uint32) (teams []*entity.Team, err error) {
	if err = gDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var teamsID []int

		if err := tx.WithContext(ctx).Table("team_user").Select("team_id").Where("user_id = ?", userID).Pluck("teams_id", &teamsID).Error; err != nil {
			return err
		}

		rows2, err2 := tx.Model(teamModel).WithContext(ctx).Where("id IN ?", teamsID).Rows()
		if err2 != nil {
			return err2
		}

		teams, err = scanTeams(tx, rows2)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return teams, nil
}

func (g *implGorm) GetTeamsByCreatorID(ctx context.Context, creatorID uint32) (teams []*entity.Team, err error) {
	rows, err := gDB.Model(teamModel).WithContext(ctx).Where("creator_id = ?", creatorID).Rows()
	if err != nil {
		return nil, err
	}

	teams, err = scanTeams(gDB, rows)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (g *implGorm) GetTeamByName(ctx context.Context, name string) (teams []*entity.Team, err error) {
	rows, err := gDB.Model(teamModel).WithContext(ctx).Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Rows()
	if err != nil {
		return nil, err
	}

	teams, err = scanTeams(gDB, rows)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (g *implGorm) UpdateTeamByID(ctx context.Context, teamID uint32, team *entity.Team) (err error) {
	if err = gDB.Model(teamModel).Where("id = ?", teamID).Updates(team).Error; err != nil {
		return err
	}

	return nil
}

const stmtAddMemberByID = "INSERT INTO team_user (team_id, user_id) VALUES (?, ?)"

func (g *implGorm) AddMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error) {
	if err = gDB.WithContext(ctx).Exec(stmtAddMemberByID, teamID, usrID).Error; err != nil {
		return err
	}

	return nil
}

const stmtRemoveMemberByID = "DELETE FROM team_user WHERE team_id = ? AND user_id = ?"

func (g *implGorm) RemoveMemberByID(ctx context.Context, teamID uint32, usrID uint32) (err error) {
	if err = gDB.WithContext(ctx).Exec(stmtRemoveMemberByID, teamID, usrID).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) DeleteTeamByID(ctx context.Context, teamID uint32) (err error) {
	if err = gDB.Model(teamModel).WithContext(ctx).Where("id = ?", teamID).Delete(entity.Team{}).Error; err != nil {
		return err
	}

	return nil
}

// PROJECT

func (g *implGorm) CreateProject(ctx context.Context, project *entity.Project) (err error) {
	if err = gDB.Model(prjModel).WithContext(ctx).Save(project).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) GetProjectByID(ctx context.Context, projectID uint32) (project *entity.Project, err error) {
	project = new(entity.Project)
	if err = gDB.Model(prjModel).WithContext(ctx).Where("id = ?", projectID).Take(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (g *implGorm) GetProjtectsByTeamID(ctx context.Context, teamID uint32) (projects []*entity.Project, err error) {
	rows, err := gDB.Model(prjModel).WithContext(ctx).Where("team_id = ?", teamID).Rows()
	if err != nil {
		return nil, err
	}

	projects, err = scanProjects(gDB, rows)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (g *implGorm) UpdateProjectByID(ctx context.Context, projectID uint32, project *entity.Project) (err error) {
	if err = gDB.Model(prjModel).WithContext(ctx).Where("id = ?", projectID).Updates(project).Error; err != nil {
		return err
	}

	return nil
}

func (g *implGorm) DeleteProjectByID(ctx context.Context, projectID uint32) (err error) {
	if err = gDB.Model(prjModel).WithContext(ctx).Where("id = ?", projectID).Delete(entity.Project{}).Error; err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- UTIL ----------------------------------------------------------

func scanUsers(gormDB *gorm.DB, rows *sql.Rows) (users []*entity.User, err error) {
	for rows.Next() {
		var user = new(entity.User)
		if err = gormDB.ScanRows(rows, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func scanTeams(gormDB *gorm.DB, rows *sql.Rows) (teams []*entity.Team, err error) {
	for rows.Next() {
		var team = new(entity.Team)
		if err = gormDB.ScanRows(rows, team); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func scanProjects(gormDB *gorm.DB, rows *sql.Rows) (projects []*entity.Project, err error) {
	for rows.Next() {
		var project = new(entity.Project)
		if err = gormDB.ScanRows(rows, project); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
