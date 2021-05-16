package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"golang.org/x/crypto/bcrypt"
)

func ginAbortWithCodeMsg(c *gin.Context, code int, message string) {
	resp := new(dto.Resp)
	resp.SetCodeMsg(code, message)
	c.AbortWithStatusJSON(code, resp)
}

func hashPassword(raw string) (hashed string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), 15)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func checkHashPassword(input string, hashed string) (ok bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
	return err == nil
}

func getClaimsFromContext(c *gin.Context) (claims *service.ClaimsDetail, err error) {
	claimsItf, ok := c.Get("Claims")
	if !ok {
		return nil, errors.New("Context claims not found")
	}

	claims, ok = claimsItf.(*service.ClaimsDetail)
	if !ok {
		return nil, errors.New("Cast claims error")
	}

	return claims, nil
}
