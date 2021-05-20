package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thanhpp/prom/cmd/portal/service"
	"github.com/thanhpp/prom/cmd/portal/webserver/dto"
	"github.com/thanhpp/prom/pkg/ccmanrpc"
	"github.com/thanhpp/prom/pkg/logger"
	"github.com/thanhpp/prom/pkg/timerpc"
)

type CardCtrl struct{}

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
		ColumnID:    req.ColumnID,
	}

	_, err = service.GetCCManSrv().CreateCard(c, int(project.ShardID), card)
	if err != nil {
		logger.Get().Errorf("Create card error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

func (cC *CardCtrl) ReorderCardInOneColumn(c *gin.Context) {
	prjID, err := getProjectIDFromParam(c)
	if err != nil {
		logger.Get().Errorf("Project ID from param error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusNotAcceptable, err.Error())
		return
	}

	req := new(dto.ReorderCardOneColumnReq)
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
		Index: strings.Trim(strings.Replace(fmt.Sprint(req.CardIndex), " ", ",", -1), "[]"),
	}

	if err = service.GetCCManSrv().UpdateColumnByID(c, int(project.ID), req.ColumnID, column); err != nil {
		logger.Get().Errorf("Update column index error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}

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
		CreatedBy:   claim.UserID,
	}
	if req.ColumnID > 0 {
		card.ColumnID = req.ColumnID
	}

	if err = service.GetCCManSrv().UpdateCardByID(c, int(project.ShardID), req.Card.ID, card); err != nil {
		logger.Get().Errorf("Update card by id error: %v", err)
		ginAbortWithCodeMsg(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := new(dto.Resp)
	resp.SetCode(http.StatusOK)
	c.JSON(http.StatusOK, resp)
}
