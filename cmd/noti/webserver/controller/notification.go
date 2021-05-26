package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thanhpp/prom/cmd/noti/repository"
	"github.com/thanhpp/prom/cmd/noti/webserver/dto"
	"github.com/thanhpp/prom/pkg/logger"
)

type NotiCtrl struct{}

func (n *NotiCtrl) GetNotiByUser(c *gin.Context) {
	userID, err := getUserIDFromQuery(c)
	if err != nil {
		logger.Get().Errorf("Get Query error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}
	page, size, err := getPageAndSizeFromQuery(c)
	if err != nil {
		logger.Get().Errorf("Get Query error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	notis, err := repository.Get().GetNotiByUserID(c, userID, page, size)
	if err != nil {
		logger.Get().Errorf("Get user noti error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	resp.SetData(gin.H{
		"notifications": notis,
	})

	c.JSON(http.StatusOK, resp)
}

func (n *NotiCtrl) GetCardNoti(c *gin.Context) {
	cardID, err := getCardIDFromQuery(c)
	if err != nil {
		logger.Get().Errorf("Get Query error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}
	page, size, err := getPageAndSizeFromQuery(c)
	if err != nil {
		logger.Get().Errorf("Get Query error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	cardNotis, err := repository.Get().GetNotiByCardID(c, cardID, page, size)
	if err != nil {
		logger.Get().Errorf("Get card noti error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	resp.SetData(gin.H{
		"notifications": cardNotis,
	})

	c.JSON(http.StatusOK, resp)
}
