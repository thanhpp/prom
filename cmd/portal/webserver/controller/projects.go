package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type ProjectCtrl struct{}

func (p *ProjectCtrl) GetAllProjectsFromTeamID(c *gin.Context) {
	teamID, err := getTeamIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("TeamID error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	projects, err := service.GetUsrManService().GetProjectsByTeamID(c, teamID)
	if err != nil {
		logger.Get().Errorf("Get projects error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK).SetData(gin.H{
		"projects": projects,
	})

	c.JSON(http.StatusOK, resp)
}

func (p *ProjectCtrl) CreateNewProject(c *gin.Context) {
	teamID, err := getTeamIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("TeamID error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.CreateProjectReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind JSON error : %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	claims, err := getClaimsFromContext(c)
	if err != nil {
		logger.Get().Errorf("Context claims error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, "Context claims error")
		return
	}

	projectID, err := service.GetUsrManService().NextProjectID(c)
	if err != nil {
		logger.Get().Errorf("Get next project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	shardID, err := service.GetCCManSrv().ChooseShardIDFromInt(int(projectID))
	if err != nil {
		logger.Get().Errorf("Generate shard id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	project := &usrmanrpc.Project{
		Name:      req.ProjectName,
		TeamID:    teamID,
		CreatorID: claims.UserID,
		ShardID:   uint32(shardID),
	}

	if err := service.GetUsrManService().NewProject(c, project); err != nil {
		logger.Get().Errorf("Create project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)

	c.JSON(http.StatusOK, resp)
}

func (p *ProjectCtrl) GetProjectDetails(c *gin.Context) {
	projectID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Get id from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	project, err := service.GetUsrManService().GetProjectByID(c, projectID)
	if err != nil {
		logger.Get().Errorf("Get project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	// choose shard
	columns, err := service.GetCCManSrv().GetColumnsByProjectID(c, int(project.ShardID), projectID)
	if err != nil {
		logger.Get().Errorf("Get columns error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK).SetData(gin.H{
		"project": project,
		"columns": columns,
	})
	c.JSON(http.StatusOK, resp)
}
