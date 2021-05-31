package dto

import (
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResp struct {
	RespError
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (r *UserLoginResp) SetToken(token string) {
	r.Data.Token = token
}

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUserNameResp struct {
	RespError
	Users []*usrmanrpc.User `json:"users"`
}

func (r *GetUserNameResp) SetData(users []*usrmanrpc.User) {
	r.Users = users
	for i := range r.Users {
		r.Users[i].HashPass = ""
	}
}
