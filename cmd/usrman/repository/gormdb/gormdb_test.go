package gormdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/thanhpp/prom/cmd/usrman/repository/entity"
	"github.com/thanhpp/prom/cmd/usrman/repository/gormdb"
)

func TestInitConnection(t *testing.T) {
	var (
		dsn = "user=thanhpp password=testthanhpp dbname=promuser host=127.0.0.1 port=5432 sslmode=disable"
	)

	if err := gormdb.GetGormDB().InitDBConnection(dsn, "INFO"); err != nil {
		t.Error(err)
		return
	}
}

func TestAutoMigrate(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx    = context.Background()
		models = []interface{}{entity.User{}, entity.Project{}, entity.Team{}}
	)

	if err := gormdb.GetGormDB().AutoMigrate(ctx, models...); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateUser(t *testing.T) {
	TestInitConnection(t)

	var (
		ctx  = context.Background()
		user = &entity.User{
			Username: "testusername3",
			HashPass: "testpass",
		}
	)

	if err := gormdb.GetGormDB().CreateUser(ctx, user); err != nil {
		t.Error(err)
		return
	}

	fmt.Println(user)
}

func TestCreateTeam(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx  = context.Background()
		team = &entity.Team{
			Name: "testteam",
		}
	)

	if err := gormdb.GetGormDB().CreateTeam(ctx, team); err != nil {
		t.Error(err)
		return
	}

	fmt.Println(team)
}

func TestAddMemberByID(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx           = context.Background()
		teamID uint32 = 1
		userID uint32 = 4
	)

	if err := gormdb.GetGormDB().AddMemberByID(ctx, teamID, userID); err != nil {
		t.Error(err)
		return
	}
}

func TestGetUsersByTeamID(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx           = context.Background()
		teamID uint32 = 1
	)

	users, err := gormdb.GetGormDB().GetUserByTeamID(ctx, teamID)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(users)
}

func TestGetTeamsByUserID(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx           = context.Background()
		userID uint32 = 1
	)

	teams, err := gormdb.GetGormDB().GetTeamsByUserID(ctx, userID)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(teams)
}

func TestRemoveMemberByID(t *testing.T) {
	TestInitConnection(t)
	var (
		ctx           = context.Background()
		teamID uint32 = 1
		userID uint32 = 4
	)

	if err := gormdb.GetGormDB().RemoveMemberByID(ctx, teamID, userID ); err != nil {
		t.Error(err)
		return
	}
}
