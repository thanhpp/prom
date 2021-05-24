package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/cmd/noti/webserver/dto"
)

func ginAbortWithCodeMsg(c *gin.Context, code int, message string) {
	resp := new(dto.Resp)
	resp.SetCodeMsg(code, message)
	c.AbortWithStatusJSON(code, resp)
}

func getUserIDFromParam(c *gin.Context) (id int, err error) {
	userIDStr := c.Param("userID")
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}

	if userIDInt == 0 {
		return 0, errors.New("Zero userID")
	}

	return userIDInt, nil
}

func getCardIDFromParam(c *gin.Context) (id int, err error) {
	cardIDStr := c.Param("cardID")
	cardIDInt, err := strconv.Atoi(cardIDStr)
	if err != nil {
		return 0, err
	}

	if cardIDInt == 0 {
		return 0, errors.New("Zero cardID")
	}

	return cardIDInt, nil
}

func getPageAndSizeFromParam(c *gin.Context) (page int, size int, err error) {
	pageStr := c.Param("page")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}

	if pageInt == 0 {
		return 0, 0, errors.New("Zero page")
	}

	sizeStr := c.Param("size")
	sizeInt, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, 0, err
	}

	if sizeInt == 0 {
		return 0, 0, errors.New("Zero size")
	}

	return pageInt, sizeInt, nil
}
