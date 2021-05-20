package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/pkg/logger"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		str := fmt.Sprintf("IP: %v | Path: %v | Method: %v | Latency: %v | Status: %v", c.ClientIP(), c.Request.URL.Path+c.Request.URL.RawQuery, c.Request.Method, time.Since(t).Seconds(), c.Writer.Status())
		if !(c.Request.URL.Path == "/health" && c.Writer.Status() == 200) {
			logger.Get().Info(str)
		}
	}
}
