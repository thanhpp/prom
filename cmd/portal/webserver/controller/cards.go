package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thanhpp/prom/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/timerpc"
)

type CardCtrl struct{}

// ------------------------------
// Create new card
// @Summary Create new card
// @Description Create new card
// @Produce json
// @Param 	Authorization	header	string					true	"jwt"
// @Param 	teamID			path	int						true	"teamID"
// @Param	projectID		path	int						true	"projectID"
// @Param	createReq		body	dto.CreateNewCardReq	true	"CreateReq"
// @Success	200	{object} dto.RespError "Create OK"
// @Tags card
// @Router /teams/:teamID/projects/:projectID/cards [POST]
func (cC *CardCtrl) CreateNewCard(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.CreateNewCardReq)
	if err = c.ShouldBindJSON(req); err != nil {
		logger.Get().Errorf("Bind JSON error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	claims, err := getClaimsFromContext(c)
	if err != nil {
		logger.Get().Errorf("Context claims error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, "Context claims error")
		return
	}

	project, err := service.GetUsrManService().GetProjectByID(c, prjID)
	if err != nil {
		logger.Get().Errorf("Get project error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	card := &ccmanrpc.Card{
		Title:       req.Card.Title,
		Description: req.Card.Description,
		AssignedTo:  req.Card.AssignedTo,
		DueDate:     timerpc.ToTimeRPC(time.Unix(int64(req.Card.DueDate), 0)),
		CreatedBy:   claims.UserID,
		ColumnID:    req.ColumnID,
	}

	id, err := service.GetCCManSrv().CreateCard(c, int(project.ShardID), card)
	if err != nil {
		logger.Get().Errorf("Create card error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	msg := &rabbitmq.NewNotificationMsg{
		CardID:  int(id),
		UserIDs: []int{int(card.AssignedTo)},
		Content: fmt.Sprintf("Card @card:%d is created by @user:%d", card.ID, card.CreatedBy),
	}

	if err := service.GetRabbitMQ().SendNewNotiMsg(msg); err != nil {
		logger.Get().Errorf("Send new noti error: %v", err)
	}

	if card.AssignedTo != 0 {
		msg := &rabbitmq.NewNotificationMsg{
			CardID:  int(id),
			UserIDs: []int{int(card.AssignedTo)},
			Content: fmt.Sprintf("Card @card:%d is assigned to @user:%d by @user:%d", card.ID, card.AssignedTo, card.CreatedBy),
		}

		if err := service.GetRabbitMQ().SendNewNotiMsg(msg); err != nil {
			logger.Get().Errorf("Send new noti error: %v", err)
		}
	}

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

// func (cC *CardCtrl) ReorderCardInOneColumn(c *gin.Context) {
// 	prjID, err := getProjectIDFromParam(c)
// 	if err != nil {
// 		logger.Get().Errorf("Project ID from param error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
// 		return
// 	}

// 	req := new(dto.ReorderCardOneColumnReq)
// 	if err = c.ShouldBindJSON(req); err != nil {
// 		logger.Get().Errorf("Bind JSON error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
// 		return
// 	}

// 	project, err := service.GetUsrManService().GetProjectByID(c, prjID)
// 	if err != nil {
// 		logger.Get().Errorf("Get project error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	column := &ccmanrpc.Column{
// 		Index: strings.Trim(strings.Replace(fmt.Sprint(req.CardIndex), " ", ",", -1), "[]"),
// 	}

// 	if err = service.GetCCManSrv().UpdateColumnByID(c, int(project.ShardID), req.ColumnID, column); err != nil {
// 		logger.Get().Errorf("Update column index error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	resp := new(dto.Resp)
// 	resp.SetCode(http.StatusOK)
// 	c.JSON(http.StatusOK, resp)
// }

// ------------------------------
// ReorderCard
// @Summary Reorder card
// @Description Reorder card if in the same column, columnID = 0
// @Produce json
// @Param 	Authorization	header	string					true	"jwt"
// @Param 	teamID			path	int						true	"teamID"
// @Param	projectID		path	int						true	"projectID"
// @Param 	reorderReq	body	dto.ReorderCard		trye	"ReorderReq"
// @Success 200 {object} dto.RespError "Reorder success"
// @Tags card
// @Router /teams/:teamID/projects/:projectID/cards/reorder [POST]
func (cC *CardCtrl) ReorderCard(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.ReorderCard)
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

	if req.ColumnID == 0 {
		if err = service.GetCCManSrv().ReorderCard(c, int(project.ShardID), req.CardID, req.AboveIdx); err != nil {
			logger.Get().Errorf("Reorcard error: %v", err)
			ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err = service.GetCCManSrv().MoveCardToCol(c, int(project.ShardID), req.CardID, req.ColumnID, req.AboveIdx); err != nil {
			logger.Get().Errorf("MoveCardToCol error: %v", err)
			ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

// // ------------------------------
// // MoveCardToCol
// func (cC *CardCtrl) MoveCardToCol(c *gin.Context) {
// 	prjID, err := getProjectIDFromParam(c)
// 	if err != nil {
// 		logger.Get().Errorf("Project ID from param error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
// 		return
// 	}

// 	claim, err := getClaimsFromContext(c)
// 	if err != nil {
// 		logger.Get().Errorf("Context claims error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
// 		return
// 	}

// 	req := new(dto.MoveCardColReq)
// 	if err = c.ShouldBindJSON(req); err != nil {
// 		logger.Get().Errorf("Bind JSON error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
// 		return
// 	}

// 	project, err := service.GetUsrManService().GetProjectByID(c, prjID)
// 	if err != nil {
// 		logger.Get().Errorf("Get project error: %v", err)
// 		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// }

// ------------------------------
// UpdateCard ...
// @Summary Update card info
// @Description Update card info (No column ID)
// @Produce json
// @Param 	Authorization	header	string					true	"jwt"
// @Param 	teamID			path	int						true	"teamID"
// @Param	projectID		path	int						true	"projectID"
// @Param 	updateReq		body	dto.UpdateCardInfoReq	true	"update info"
// @Success 200 {object} dto.RespError "Update OK"
// @Tags card
// @Router /teams/:teamID/projects/:projectID/cards [PATCH]
func (cC *CardCtrl) UpdateCard(c *gin.Context) {
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

	req := new(dto.UpdateCardInfoReq)
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

	card := &ccmanrpc.Card{
		Title:       req.Card.Title,
		Description: req.Card.Description,
		AssignedTo:  req.Card.AssignedTo,
		DueDate:     timerpc.ToTimeRPC(time.Unix(int64(req.Card.DueDate), 0)),
	}

	// prevent column update
	card.ColumnID = 0

	if err = service.GetCCManSrv().UpdateCardByID(c, int(project.ShardID), req.Card.ID, card); err != nil {
		logger.Get().Errorf("Update card by id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	msg := &rabbitmq.NewNotificationMsg{
		CardID:  int(card.ID),
		Content: fmt.Sprintf("Card @card:%d is updated by @user:%d", req.Card.ID, claim.UserID),
	}

	if err := service.GetRabbitMQ().SendNewNotiMsg(msg); err != nil {
		logger.Get().Errorf("Send new noti error: %v", err)
	}

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

// ------------------------------
// DeleteCardByID
// @Summary Delete card by ID
// @Description Delete card by ID
// @Produce json
// @Param 	Authorization	header	string					true	"jwt"
// @Param 	teamID			path	int						true	"teamID"
// @Param	projectID		path	int						true	"projectID"
// @Param 	deleteReq		body 	dto.DeleteCardReq		true	"delete request"
// @Success 200 {object} dto.RespError "delete success"
// @Tags card
// @Router /teams/:teamID/projects/:projectID/cards [DELETE]
func (cC *CardCtrl) DeleteCard(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	var req = new(dto.DeleteCardReq)
	if err := c.ShouldBindJSON(req); err != nil {
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

	if err = service.GetCCManSrv().DeleteCardByID(c, int(project.ShardID), req.CardID); err != nil {
		logger.Get().Errorf("Delete card error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.RespError)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}
