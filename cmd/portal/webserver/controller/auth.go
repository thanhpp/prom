package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thanhpp/prom/cmd/portal/repository"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/logger"
)

type AuthCtrl struct{}

func (a AuthCtrl) Login(c *gin.Context) {
	req := new(dto.UserLoginReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind JSON error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, "Invalid request")
		return
	}

	// get user
	usr, err := service.GetUsrManService().Login(c, req.Username, "%")
	if err != nil {
		logger.Get().Errorf("Usrman service error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !checkHashPassword(req.Password, usr.HashPass) {
		logger.Get().Errorf("Check password error: %s. Expected: %s", req.Password, usr.HashPass)
		ginAbortWithCodeMsg(c, http.StatusUnauthorized, "Wrong password")
		return
	}

	// create token
	var (
		now = time.Now()
		dur = time.Minute * 15
		exp = now.Add(dur)
	)
	claims := &service.ClaimsDetail{
		UUID:   uuid.New().String(),
		UserID: usr.ID,
		// IssuedAt:  now.UnixNano(),
		ExpiredAt: exp.UnixNano(),
	}
	token, err := service.GetJWTSrv().CreateToken(claims)
	if err != nil {
		logger.Get().Errorf("Create token error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	// save to redis
	if err := repository.GetRedis().SetKey(c, claims.UUID, "OK", dur); err != nil {
		logger.Get().Errorf("Set redis key error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusUnauthorized, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCodeMsg(http.StatusOK, "")
	resp.Data = gin.H{
		"token": token,
	}

	c.JSON(http.StatusOK, resp)
}

func (a AuthCtrl) Logout(c *gin.Context) {
	fmt.Println("Checkpoint")

	val, ok := c.Get("Claims")
	if !ok {
		logger.Get().Error("Get claims from context error")
		ginAbortWithCodeMsg(c, http.StatusUnauthorized, "Claims not found")
		return
	}
	claims, ok := val.(*service.ClaimsDetail)
	if !ok {
		logger.Get().Errorf("Cast context claims error: %v", claims)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, "Cast context claims error")
		return
	}

	if err := repository.GetRedis().DeleteKey(c, claims.UUID); err != nil {
		logger.Get().Errorf("Delete key error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
