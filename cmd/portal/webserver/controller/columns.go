package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/logger"
)

type ColumnCtrl struct{}

// ------------------------------
// CreateNewColumn ...
// @Summary Create new column
// @Description Create new column
// @Produce json
// @Param 	Authorization	header	string					true	"jwt"
// @Param 	teamID			path	int						true	"teamID"
// @Param	projectID		path	int						true	"projectID"
// @Param 	createReq		body	dto.CreateNewColumnReq	true	"CreateReq"
// @Success 200 {object} dto.RespError "Create success"
// @Tags column
// @Router /teams/:teamID/projects/:projectID/columns [POST]
func (cC *ColumnCtrl) CreateNewColumn(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	claim, err := getClaimsFromContext(c)
	if err != nil {
		logger.Get().Errorf("Context claims error: %v", err)
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
		ProjectID: prjID,
		Title:     req.ColumnName,
		CreatedBy: claim.UserID,
	}
	_, err = service.GetCCManSrv().CreateColumn(c, int(project.ShardID), column)
	if err != nil {
		logger.Get().Errorf("Create column error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	// // update project index
	// project.Index += fmt.Sprintf("%d,", id)
	// if err = service.GetUsrManService().UpdateProject(c, project); err != nil {
	// 	logger.Get().Errorf("Update project index error: %v", err)
	// 	ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

// ------------------------------
// ReorderColumn ...
// @Summary Reorder 1 column
// @Description Reorder 1 column
// @Produce json
// @Param 	Authorization	header	string					true	"jwt"
// @Param 	teamID			path	int						true	"teamID"
// @Param	projectID		path	int						true	"projectID"
// @Param 	reorderReq		body	dto.UpdateColumnIndex	true	"reorderReq"
// @Success 200 {object} dto.RespError "ReorderOK"
// @Tags column
// @Router /teams/:teamID/projects/:projectID/columns/reorder  [POST]
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

	project, err := service.GetUsrManService().GetProjectByID(c, prjID)
	if err != nil {
		logger.Get().Errorf("Get project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := service.GetCCManSrv().ReorderColumn(c, int(project.ShardID), req.ColumnID, req.NextOfIdx); err != nil {
		logger.Get().Errorf("Update column index error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

// ------------------------------
// UpdateColumn ...
// @Summary Update column by id
// @Description Update column by id
// @Produce json
// @Param 	Authorization	header	string				true	"jwt"
// @Param 	teamID			path	int					true	"teamID"
// @Param	projectID		path	int					true	"projectID"
// @Param 	updateReq 		body 	dto.UpdateColumnReq true "update column info"
// @Success 200 {object} dto.RespError "Update OK"
// @Tags column
// @Router /teams/:teamID/projects/:projectID/columns  [PATCH]
func (cC *ColumnCtrl) UpdateColumn(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.UpdateColumnReq)
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

	if err = service.GetCCManSrv().UpdateColumnByID(c, int(project.ShardID), req.ColumnID, req.Column); err != nil {
		logger.Get().Errorf("UpdateColumnByID error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

// ------------------------------
// DeleteColumn ...
// @Summary Delete column by id
// @Description Delete column by id and all card in column
// @Produce json
// @Param 	Authorization	header	string				true	"jwt"
// @Param 	teamID			path	int					true	"teamID"
// @Param	projectID		path	int					true	"projectID"
// @Param 	deleteReq		body	dto.DeleteColumn	true	"deleteReq"
// @Success 200 {object} dto.RespError
// @Tags column
// @Router /teams/:teamID/projects/:projectID/columns [DELETE]
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

	if err = service.GetCCManSrv().DeleteColumnByID(c, int(project.ShardID), req.ColumnID); err != nil {
		logger.Get().Errorf("Delete column error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	// // update project index
	// project.Index = strings.ReplaceAll(project.Index, fmt.Sprintf("%d,", req.ColumnID), "")
	// fmt.Println(project.Index)
	// if err = service.GetUsrManService().UpdateProject(c, project); err != nil {
	// 	logger.Get().Errorf("Update project index error: %v", err)
	// 	ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}
