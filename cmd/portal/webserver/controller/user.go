package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/logger"
)

type UserCtrl struct{}

func (u UserCtrl) CreateNewUser(c *gin.Context) {
	req := new(dto.CreateUserReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind JSON error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	pwd, err := hashPassword(req.Password)
	if err != nil {
		logger.Get().Errorf("Hash password error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := service.GetUsrManService().NewUser(c, req.Username, pwd); err != nil {
		logger.Get().Errorf("Create new user error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

func (u UserCtrl) UpdateUser(c *gin.Context) {
	req := new(dto.UpdateUserReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind JSON error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	claims, err := getClaimsFromContext(c)
	if err != nil {
		logger.Get().Errorf("Get claims error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := service.GetUsrManService().UpdateUsername(c, claims.UserID, req.Username); err != nil {
		logger.Get().Errorf("Update username error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	pwd, err := hashPassword(req.Password)
	if err != nil {
		logger.Get().Errorf("Hash password error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := service.GetUsrManService().UpdatePassword(c, claims.UserID, pwd); err != nil {
		logger.Get().Errorf("Update password error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

func (u UserCtrl) GetUserName(c *gin.Context) {
	usrname := c.Query("username")
	if len(usrname) == 0 {
		logger.Get().Error("Empty username param")
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, "Empty username param")
		return
	}

	users, err := service.GetUsrManService().GetUsersByPattern(c, usrname)
	if err != nil {
		logger.Get().Errorf("Get pattern error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	resp.SetData(gin.H{
		"users": users,
	})
	c.JSON(http.StatusOK, resp)
}
