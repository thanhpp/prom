package service

import (
	"context"
	"log"
	"reflect"

	"github.com/thanhpp/prom/pkg/errconst"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
	"google.golang.org/grpc"
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

func (us usrManSrv) respError(err error, resp interface{}) error {
	if resp == nil {
		return errconst.ServiceError{Srv: us.name(), Err: err}
	}

	defer func() {
		v := recover()
		log.Printf("Recover from usrManSrv.respErr: %v \n", v)
	}()

	val := reflect.ValueOf(resp)
	code := val.FieldByName("Code").Int()
	msg := val.FieldByName("Message").String()

	return errconst.ServiceError{Srv: us.name(), Err: err, Code: int32(code), Msg: msg}
}

type iUsrManSrv interface {
	// user
	Login(ctx context.Context, username string, pass string) (user *usrmanrpc.User, err error)
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

var implUsrManSrv = new(usrManSrv)

func SetUsrManService(ctx context.Context, target string) (err error) {
	conn, err := grpc.DialContext(ctx, target, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
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
	if err != nil || resp.Code != errconst.RPCSuccessCode {
		return nil, uS.respError(err, resp)
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

	resp, err := uS.client().CreateUser(ctx, in)
	if err != nil || resp.Code != errconst.RPCSuccessCode {
		return uS.respError(err, resp)
	}

	return nil
}

func (uS *usrManSrv) UpdateUsername(ctx context.Context, userID uint32, username string) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) UpdatePassword(ctx context.Context, userID uint32, password string) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) GetTeamsByUserID(ctx context.Context, userID uint32) (teams []*usrmanrpc.Team, err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) GetTeamMembersByID(ctx context.Context, teamID uint32) (users []*usrmanrpc.User, err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) AddMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) RemoveMemberByID(ctx context.Context, teamID uint32, userID uint32) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) GetProjectsByTeamID(ctx context.Context, teamID uint32) (projects []*usrmanrpc.Project, err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) NewProject(ctx context.Context, project *usrmanrpc.Project) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) UpdateProject(ctx context.Context, project *usrmanrpc.Project) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}

func (uS *usrManSrv) DeleteProjectByID(ctx context.Context, projectID uint32) (err error) {
	if ctx.Err() != nil {
		return
	}

	return
}
