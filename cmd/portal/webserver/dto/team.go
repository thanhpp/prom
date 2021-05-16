package dto

type CreateNewTeamReq struct {
	TeamName string `json:"teamName"`
}

type UpdateTeamMemberReq struct {
	Op       string `json:"op"`
	MemberID uint32 `json:"memberID"`
}
