package controller

import (
	"net/http"

	"github.com/thanhpp/prom/pkg/usrmanrpc"

	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/pkg/logger"
)

type TeamCtrl struct{}

func (t *TeamCtrl) GetAllTeamByUserID(c *gin.Context) {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		logger.Get().Errorf("Context claims error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	teams, err := service.GetUsrManService().GetTeamsByUserID(c, claims.UserID)
	if err != nil {
		logger.Get().Errorf("Get teams error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK).SetData(gin.H{
		"teams": teams,
	})

	c.JSON(http.StatusOK, resp)
}

func (t *TeamCtrl) CreateNewTeam(c *gin.Context) {
	req := new(dto.CreateNewTeamReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind json error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	claims, err := getClaimsFromContext(c)
	if err != nil {
		logger.Get().Errorf("Context claims error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	team := &usrmanrpc.Team{
		Name:      req.TeamName,
		CreatorID: claims.UserID,
	}

	if err := service.GetUsrManService().CreateNewTeam(c, team); err != nil {
		logger.Get().Errorf("Create team error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)

	c.JSON(http.StatusOK, resp)
}

func (t *TeamCtrl) GetTeamByID(c *gin.Context) {
	teamID, err := getTeamIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Team id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	team, err := service.GetUsrManService().GetTeamByID(c, teamID)
	if err != nil {
		logger.Get().Errorf("Get team by ID error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK).SetData(gin.H{
		"team": team,
	})
	c.JSON(http.StatusOK, resp)
}

func (t *TeamCtrl) EditMember(c *gin.Context) {
	teamID, err := getTeamIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Team id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.UpdateTeamMemberReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind json error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	switch req.Op {
	case "add":
		if err := service.GetUsrManService().AddMemberByID(c, teamID, req.MemberID); err != nil {
			logger.Get().Errorf("Add member error: %v", err)
			ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
			return
		}
	case "remove":
		if err := service.GetUsrManService().RemoveMemberByID(c, teamID, req.MemberID); err != nil {
			logger.Get().Errorf("Remove member error: %v", err)
			ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		logger.Get().Errorf("Invalid option: %s", req.Op)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, "Invalid option")
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

func (t *TeamCtrl) DeleteTeam(c *gin.Context) {
	teamID, err := getTeamIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Team id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	if err := service.GetUsrManService().DeleteTeamByID(c, teamID); err != nil {
		logger.Get().Errorf("Delete team by id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}
