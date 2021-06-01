package service

import (
	"context"
	"time"

	"github.com/thanhpp/prom/pkg/etcdclient"

	"google.golang.org/grpc"

	"github.com/thanhpp/prom/pkg/errconst"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

// ----------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- USER MANAGER SERVICE ----------------------------------------------------------

type usrManSrv struct {
	grpcClient usrmanrpc.UsrManSrvClient
}

func (uS usrManSrv) client() (client usrmanrpc.UsrManSrvClient) {
	return uS.grpcClient
}

func (us usrManSrv) name() string {
	return "User manager service"
}

func (us usrManSrv) respError(err error, respCode int32, respMessage string) error {
	return errconst.ServiceError{Srv: us.name(), Err: err, Code: respCode, Msg: respMessage}
}

func (us usrManSrv) error(err error) error {
	return errconst.ServiceError{Srv: us.name(), Err: err, Msg: err.Error()}
}

type iUsrManSrv interface {
	// user
	Login(ctx context.Context, username string, pass string) (user *usrmanrpc.User, err error)
	GetUsersByPattern(ctx context.Context, pattern string) (users []*usrmanrpc.User, err error)
	NewUser(ctx context.Context, username string, pass string) (err error)
	UpdateUsername(ctx context.Context, userID uint32, username string) (err error)
	UpdatePassword(ctx context.Context, userID uint32, password string) (err error)

	// team
	CreateNewTeam(ctx context.Context, team *usrmanrpc.Team) (err error)
	GetTeamByID(ctx context.Context, teamID uint32) (team *usrmanrpc.Team, err error)
	GetTeamsByUserID(ctx context.Context, userID uint32) (teams []*usrmanrpc.Team, err error)
	GetTeamMembersByID(ctx context.Context, teamID uint32) (users []*usrmanrpc.User, err error)
	AddMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error)
	RemoveMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error)
	DeleteTeamByID(ctx context.Context, teamID uint32) (err error)

	// project
	GetProjectByID(ctx context.Context, projectID uint32) (project *usrmanrpc.Project, err error)
	GetRecentCreatedProjectByUserID(ctx context.Context, userID uint32, recent uint) (projects []*usrmanrpc.Project, err error)
	GetProjectsByTeamID(ctx context.Context, teamID uint32) (projects []*usrmanrpc.Project, err error)
	NextProjectID(ctx context.Context) (id uint32, err error)
	NewProject(ctx context.Context, project *usrmanrpc.Project) (err error)
	UpdateProject(ctx context.Context, project *usrmanrpc.Project) (err error)
	ReorderProjectColumns(ctx context.Context, projectID uint32, columnsIdx string) (err error)
	AddColumnsToProject(ctx context.Context, projectID uint32, columnID uint32) (err error)
	DeleteProjectByID(ctx context.Context, projectID uint32) (err error)
}

var implUsrManSrv = new(usrManSrv)

func SetUsrManService(ctx context.Context, target string) (err error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	b := grpc.RoundRobin(etcdclient.Get().Resolver())
	conn, err := grpc.DialContext(newCtx, target, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		cancel()
		return err
	}

	implUsrManSrv.grpcClient = usrmanrpc.NewUsrManSrvClient(conn)
	return nil
}

func GetUsrManService() iUsrManSrv {
	return implUsrManSrv
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNCTIONS ----------------------------------------------------------

func (uS *usrManSrv) Login(ctx context.Context, username string, pass string) (user *usrmanrpc.User, err error) {
	if ctx.Err() != nil {
		return nil, errconst.ServiceError{Srv: uS.name(), Err: ctx.Err(), Msg: "Context error"}
	}

	in := &usrmanrpc.GetUserByUsernamePassReq{
		Username: username,
		Pass:     pass,
	}

	resp, err := uS.client().GetUserByUsernamePass(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	if resp.User == nil {
		return nil, uS.respError(err, resp.Code, resp.Message)
	}

	return resp.User, nil
}

func (uS *usrManSrv) NewUser(ctx context.Context, username string, pass string) (err error) {
	if ctx.Err() != nil {
		return errconst.ServiceError{Srv: uS.name(), Err: ctx.Err(), Msg: "Context error"}
	}

	in := &usrmanrpc.CreateUserReq{
		User: &usrmanrpc.User{
			Username: username,
			HashPass: pass,
		},
	}

	_, err = uS.client().CreateUser(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}

func (uS *usrManSrv) GetUsersByPattern(ctx context.Context, pattern string) (users []*usrmanrpc.User, err error) {
	if ctx.Err() != nil {
		return nil, errconst.ServiceError{Srv: uS.name(), Err: ctx.Err(), Msg: "Context error"}
	}

	in := &usrmanrpc.GetUserByPatternReq{
		Pattern: pattern,
	}

	resp, err := uS.client().GetUserByPattern(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Users, nil
}

func (uS *usrManSrv) UpdateUsername(ctx context.Context, userID uint32, username string) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.UpdateUserByIDReq{
		UserID: userID,
		User: &usrmanrpc.User{
			Username: username,
		},
	}

	_, err = uS.client().UpdateUserByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}

func (uS *usrManSrv) UpdatePassword(ctx context.Context, userID uint32, password string) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.UpdateUserByIDReq{
		UserID: userID,
		User: &usrmanrpc.User{
			HashPass: password,
		},
	}

	_, err = uS.client().UpdateUserByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return
}

// ------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- TEAM ----------------------------------------------------------

func (uS *usrManSrv) CreateNewTeam(ctx context.Context, team *usrmanrpc.Team) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.CreateTeamReq{
		Team: team,
	}

	_, err = uS.client().CreateTeam(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}

func (uS *usrManSrv) GetTeamByID(ctx context.Context, teamID uint32) (team *usrmanrpc.Team, err error) {
	if ctx.Err() != nil {
		return nil, uS.error(ctx.Err())
	}

	in := &usrmanrpc.GetTeamByIDReq{
		TeamID: teamID,
	}

	resp, err := uS.client().GetTeamByID(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Team, nil
}

func (uS *usrManSrv) GetTeamsByUserID(ctx context.Context, userID uint32) (teams []*usrmanrpc.Team, err error) {
	if ctx.Err() != nil {
		return nil, uS.error(ctx.Err())
	}

	in := &usrmanrpc.GetTeamsByUserIDReq{
		UserID: userID,
	}

	resp, err := uS.client().GetTeamsByUserID(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Teams, nil
}

func (uS *usrManSrv) GetTeamMembersByID(ctx context.Context, teamID uint32) (users []*usrmanrpc.User, err error) {
	if ctx.Err() != nil {
		return nil, uS.error(ctx.Err())
	}

	in := &usrmanrpc.GetUserByTeamIDReq{
		TeamID: teamID,
	}

	resp, err := uS.client().GetUserByTeamID(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Users, nil
}

func (uS *usrManSrv) AddMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.AddMemberByIDReq{
		TeamID: teamID,
		UserID: userID,
	}

	_, err = uS.client().AddMemberByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}

func (uS *usrManSrv) RemoveMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.RemoveMemberByIDReq{
		TeamID: teamID,
		UserID: userID,
	}

	_, err = uS.client().RemoveMemberByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return
}

func (uS *usrManSrv) DeleteTeamByID(ctx context.Context, teamID uint32) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.DeleteTeamByIDReq{
		TeamID: teamID,
	}

	_, err = uS.client().DeleteTeamByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return
}

// ----------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- PROJECT ----------------------------------------------------------

func (uS *usrManSrv) GetProjectByID(ctx context.Context, projectID uint32) (project *usrmanrpc.Project, err error) {
	if ctx.Err() != nil {
		return nil, uS.error(ctx.Err())
	}

	in := &usrmanrpc.GetProjectByIDReq{
		ProjectID: projectID,
	}

	resp, err := uS.client().GetProjectByID(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Project, nil
}

func (uS *usrManSrv) GetRecentCreatedProjectByUserID(ctx context.Context, userID uint32, recent uint) (projects []*usrmanrpc.Project, err error) {
	if ctx.Err() != nil {
		return nil, uS.error(ctx.Err())
	}

	in := &usrmanrpc.GetRecentCreatedProjectByUserIDReq{
		UserID: userID,
		Recent: uint32(recent),
	}

	resp, err := uS.client().GetRecentCreatedProjectByUserID(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Projects, nil
}

func (uS *usrManSrv) GetProjectsByTeamID(ctx context.Context, teamID uint32) (projects []*usrmanrpc.Project, err error) {
	if ctx.Err() != nil {
		return nil, uS.error(ctx.Err())
	}

	in := &usrmanrpc.GetProjtectsByTeamIDReq{
		TeamID: teamID,
	}

	resp, err := uS.client().GetProjtectsByTeamID(ctx, in)
	if err != nil {
		return nil, uS.error(err)
	}

	return resp.Projects, nil
}

func (uS *usrManSrv) NextProjectID(ctx context.Context) (id uint32, err error) {
	if ctx.Err() != nil {
		return 0, uS.error(ctx.Err())
	}

	in := &usrmanrpc.NextProjectIDReq{}

	resp, err := uS.client().NextProjectID(ctx, in)
	if err != nil {
		return 0, uS.error(err)
	}

	return uint32(resp.NextID), nil
}

func (uS *usrManSrv) NewProject(ctx context.Context, project *usrmanrpc.Project) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.CreateProjectReq{
		Project: project,
	}

	_, err = uS.client().CreateProject(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}

func (uS *usrManSrv) UpdateProject(ctx context.Context, project *usrmanrpc.Project) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.UpdateProjectByIDReq{
		ProjectID: project.ID,
		Project:   project,
	}

	_, err = uS.client().UpdateProjectByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}

func (uS *usrManSrv) ReorderProjectColumns(ctx context.Context, projectID uint32, columnsIdx string) (err error) {

	return nil
}

func (uS *usrManSrv) AddColumnsToProject(ctx context.Context, projectID uint32, columnID uint32) (err error) {

	return nil
}

func (uS *usrManSrv) DeleteProjectByID(ctx context.Context, projectID uint32) (err error) {
	if ctx.Err() != nil {
		return uS.error(ctx.Err())
	}

	in := &usrmanrpc.DeleteProjectByIDReq{
		ProjectID: projectID,
	}

	_, err = uS.client().DeleteProjectByID(ctx, in)
	if err != nil {
		return uS.error(err)
	}

	return nil
}
