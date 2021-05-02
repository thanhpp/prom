package service

type UsrManSrv struct {
}

type iUsrManSrv interface {
	// user
	Login(username string, pass string) (jwt string, err error)
	Logout(userID uint32) (err error)
	NewUser(username string, pass string) (err error)
	UpdateUsername(userID uint32, username string) (err error)
	UpdatePassword(userID uint32, password string) (err error)

	// team
	GetTeamsByUsersID(userID uint32)
}
