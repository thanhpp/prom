package dto

import (
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type CreateProjectReq struct {
	ProjectName string `json:"projectName"`
}

type GetAllProjectFromTeamIDResp struct {
	RespError
	Projects []*usrmanrpc.Project `json:"projects"`
}

func (r *GetAllProjectFromTeamIDResp) SetData(prjs []*usrmanrpc.Project) {
	r.Projects = prjs
}

type GetProjectDetailsResp struct {
	RespError
	Project *usrmanrpc.Project `json:"project"`
	Column  []*ccmanrpc.Column `json:"columns"`
}

func (r *GetProjectDetailsResp) SetData(prj *usrmanrpc.Project, cols []*ccmanrpc.Column) {
	r.Project = prj
	r.Column = cols
}

type GetRecentCreatedProjectByUserIDResp struct {
	RespError
	Projects []*usrmanrpc.Project `json:"projects"`
}

func (r *GetRecentCreatedProjectByUserIDResp) SetData(prjs []*usrmanrpc.Project) {
	r.Projects = prjs
}
