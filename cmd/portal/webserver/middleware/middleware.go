package middleware

import (
	"net/http"
	"strings"

	"github.com/thanhpp/prom/pkg/logger"

	"github.com/gin-gonic/gin"

	"github.com/thanhpp/prom/cmd/portal/repository"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
)

const ()

type AuthMw struct{}

func (a AuthMw) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract token
		req := c.Request
		token := req.Header.Get("Authorization")
		if i := strings.Index(token, "Bearer "); i == -1 {
			logger.Get().Errorf("Invalid token: %s", token)
			resp := new(dto.Resp)
			resp.SetCodeMsg(http.StatusUnauthorized, "Invalid Authorization Header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		// claims validate with key
		claims, err := service.GetJWTSrv().GetClaimsValidate(token[7:])
		if err != nil {
			logger.Get().Errorf("Validate token error: %v", err)
			resp := new(dto.Resp)
			resp.SetCodeMsg(http.StatusUnauthorized, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		// redis logout validate
		_, err = repository.GetRedis().GetValue(c, claims.UUID)
		if err != nil {
			logger.Get().Errorf("Redis get value error: %v", err)
			resp := new(dto.Resp)
			resp.SetCodeMsg(http.StatusUnauthorized, "Expired Token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		c.Set("Claims", claims)

		c.Next()
	}
}
