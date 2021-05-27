package router

import (
	"time"

	"github.com/thanhpp/prom/cmd/portal/webserver/controller"
	"github.com/thanhpp/prom/cmd/portal/webserver/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() (routers *gin.Engine) {
	routers = gin.New()
	routers.Use(gin.Recovery())
	routers.Use(middleware.LogRequest())

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
	authCtrl := new(controller.AuthCtrl)
	authMw := new(middleware.AuthMw)

	routers.POST("/login", authCtrl.Login)
	routers.GET("/logout", authMw.ValidateToken(), authCtrl.Logout) // NOTE: JWT VALIDATE

	// user
	usrCtrl := new(controller.UserCtrl)
	user := routers.Group("/user")
	{
		user.POST("", usrCtrl.CreateNewUser)
		user.Use(authMw.ValidateToken())
		user.PATCH("", usrCtrl.UpdateUser)
		user.GET("", usrCtrl.GetUserName) // NOTE: query username
	}

	// JWT from here
	routers.Use(authMw.ValidateToken())
	teamCtrl := new(controller.TeamCtrl)
	teams := routers.Group("/teams")
	{
		teams.GET("", teamCtrl.GetAllTeamByUserID)
		teams.POST("", teamCtrl.CreateNewTeam)
		teamsID := teams.Group("/:teamID")
		{
			teamsID.GET("", teamCtrl.GetTeamByID)
			teamsID.PUT("", teamCtrl.EditMember)
			teamsID.DELETE("", teamCtrl.DeleteTeam)

			prjCtrl := new(controller.ProjectCtrl)
			projects := teamsID.Group("/projects")
			{
				projects.GET("", prjCtrl.GetAllProjectsFromTeamID)
				projects.POST("", prjCtrl.CreateNewProject)

				projectID := projects.Group("/:projectID")
				{
					projectID.GET("", prjCtrl.GetProjectDetails)

					colCtrl := new(controller.ColumnCtrl)
					columns := projectID.Group("/columns")
					{
						columns.POST("", colCtrl.CreateNewColumn)
						columns.POST("/reorder", colCtrl.ReorderColumns)
						columns.DELETE("", colCtrl.DeleteColumn)
					}

					cardCtrl := new(controller.CardCtrl)
					cards := projectID.Group("/cards")
					{
						cards.POST("", cardCtrl.CreateNewCard)
						cards.POST("/reorder", cardCtrl.ReorderCard)
						cards.PATCH("", cardCtrl.UpdateCard)
					}
				}
			}
		}

	}
	return routers
}
