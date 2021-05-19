package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/usrmanrpc"
)

type ColumnCtrl struct{}

func (cC *ColumnCtrl) CreateNewColumn(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.CreateNewColumnReq)
	if err = c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind JSON error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	project, err := service.GetUsrManService().GetProjectByID(c, prjID)
	if err != nil {
		logger.Get().Errorf("Get project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	column := &ccmanrpc.Column{
		Title: req.ColumnName,
	}
	if err = service.GetCCManSrv().CreateColumn(c, int(project.ShardID), column); err != nil {
		logger.Get().Errorf("Create column error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

func (cC *ColumnCtrl) ReorderColumns(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.UpdateColumnIndex)
	if err = c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind json error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	project := &usrmanrpc.Project{
		ID:    prjID,
		Index: strings.Trim(strings.Replace(fmt.Sprint(req.ColumnIndex), " ", ",", -1), "[]"),
	}
	if err = service.GetUsrManService().UpdateProject(c, project); err != nil {
		logger.Get().Errorf("Update project index error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

func (cC *ColumnCtrl) DeleteColumn(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.DeleteColumn)
	if err = c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind json error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	project, err := service.GetUsrManService().GetProjectByID(c, prjID)
	if err != nil {
		logger.Get().Errorf("Get project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	if req.MoveTo == 0 {
		if err = service.GetCCManSrv().DeleteColumnByID(c, int(project.ShardID), req.ColumnID); err != nil {
			logger.Get().Errorf("Delete column error: %v", err)
			ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err = service.GetCCManSrv().DeleteColumnByIDAndMove(c, int(project.ShardID), req.ColumnID, req.MoveTo); err != nil {
			logger.Get().Errorf("Delete column error: %v", err)
			ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}
