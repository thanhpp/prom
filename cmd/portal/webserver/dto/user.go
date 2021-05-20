package dto

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
