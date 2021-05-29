package dto

import (
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type CreateNewTeamReq struct {
	TeamName string `json:"teamName"`
}

type UpdateTeamMemberReq struct {
	Op       string `json:"op"`
	MemberID uint32 `json:"memberID"`
}

type GetAllTeamByUserIDResp struct {
	RespError
	Teams []*usrmanrpc.Team `json:"teams"`
}

func (r *GetAllTeamByUserIDResp) SetData(teams []*usrmanrpc.Team) {
	r.Teams = teams
}

type GetTeamByIDResp struct {
	RespError
	Team *usrmanrpc.Team `json:"team"`
}

func (r *GetTeamByIDResp) SetData(team *usrmanrpc.Team) {
	r.Team = team
}
