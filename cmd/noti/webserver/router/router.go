package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/thanhpp/prom/cmd/noti/webserver/controller"
)

func NewRouter() (routers *gin.Engine) {
	routers = gin.New()
	//CORS
	routers.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin,DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	notiCtrl := new(controller.NotiCtrl)

	notiGr := routers.Group("notifications")
	{
		notiGr.GET("/user", notiCtrl.GetNotiByUser)
		notiGr.GET("/card", notiCtrl.GetCardNoti)
	}

	return routers
}
