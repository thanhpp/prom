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

func getUserIDFromQuery(c *gin.Context) (id int, err error) {
	userIDStr := c.Query("userID")
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}

	if userIDInt == 0 {
		return 0, errors.New("Zero userID")
	}

	return userIDInt, nil
}

func getCardIDFromQuery(c *gin.Context) (id int, err error) {
	cardIDStr := c.Query("cardID")
	cardIDInt, err := strconv.Atoi(cardIDStr)
	if err != nil {
		return 0, err
	}

	if cardIDInt == 0 {
		return 0, errors.New("Zero cardID")
	}

	return cardIDInt, nil
}

func getPageAndSizeFromQuery(c *gin.Context) (page int, size int, err error) {
	page = 1
	pageStr := c.Query("page")
	if len(pageStr) > 0 {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, err
		}
		if page <= 0 {
			return 0, 0, errors.New("Zero page")
		}
	}

	size = 1
	sizeStr := c.Query("size")
	if len(sizeStr) > 0 {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return 0, 0, err
		}

		if size <= 0 {
			return 0, 0, errors.New("Zero size")
		}
	}

	return page, size, nil
}
