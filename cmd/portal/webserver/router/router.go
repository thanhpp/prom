package router

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/thanhpp/prom/pkg/logger"
)

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		str := fmt.Sprintf("IP: %v | Path: %v | Method: %v | Latency: %v | Status: %v", c.ClientIP(), c.Request.URL.Path+c.Request.URL.RawQuery, c.Request.Method, time.Since(t).Seconds(), c.Writer.Status())
		if !(c.Request.URL.Path == "/health" && c.Writer.Status() == 200) {
			logger.Get().Info(str)
		}
	}
}

func NewRouter() (routers *gin.Engine) {
	routers = gin.New()
	routers.Use(gin.Recovery())
	routers.Use(logRequest())

	//CORS
	routers.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin,DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// login/logout
	routers.POST("/login")
	routers.POST("/logout") // NOTE: JWT VALIDATE

	// user
	user := routers.Group("/user")
	{
		user.POST("")
		// JWT VALIDATE
		user.PATCH("")
		user.GET("") // NOTE: query username
	}

	// JWT from here
	teams := routers.Group("/teams")
	{
		teams.GET("")
		teams.POST("")
		teamsID := teams.Group("/:teamID")
		{
			teamsID.GET("/:teamID")
			teamsID.PUT("/:teamID")
			teamsID.DELETE("/:teamID")

			projects := teamsID.Group("/projects")
			{
				projects.GET("")
				projects.POST("")

				projectID := projects.Group("/:projectID")
				{
					projectID.GET("")
					projectID.POST("")

					columns := projectID.Group("/columns")
					{
						columns.POST("")
						columns.PATCH("")
						columns.DELETE("")
					}

					cards := projectID.Group("/cards")
					{
						cards.POST("")
						cards.PUT("")
						cards.PATCH("")
					}
				}
			}
		}

	}
	return routers
}
