package controller

import (
	"event-management-system/middleware"
	"event-management-system/models"
	"event-management-system/usecase"
	modelutil "event-management-system/utils/model_util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketUseCase  usecase.TicketUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewTicketController(ticketUseCase usecase.TicketUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *TicketController {
	return &TicketController{ticketUseCase: ticketUseCase, rg: rg, authMiddleware: authMiddleware}
}

func (tc *TicketController) Route() {
	tc.rg.POST("/tickets", tc.authMiddleware.RequireToken("admin"), tc.createTicket)
}

func (tc *TicketController) createTicket(ctx *gin.Context) {
	var payload []models.Tickets

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	ticket, err := tc.ticketUseCase.CreateTicket(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ticketResponse := make([]gin.H, 0)

	for _, t := range ticket {
		ticketResponse = append(ticketResponse, gin.H{
			"id":         t.Id,
			"ticketUuid": t.TikcetUuid,
			"ticketType": t.TicketType,
			"price":      t.Price,
			"quota":      t.Quota,
			"status":     t.Status,
			"createdAt":  t.CreatedAt,
			"updatedAt":  nil,
			"eventId":    t.EventID,
		})
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success create ticket", ticketResponse, true))
}
