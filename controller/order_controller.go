package controller

import (
	"event-management-system/middleware"
	"event-management-system/models"
	"event-management-system/usecase"
	modelutil "event-management-system/utils/model_util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderUseCase   usecase.OrderUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewOrderController(orderUseCase usecase.OrderUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *OrderController {
	return &OrderController{orderUseCase: orderUseCase, rg: rg, authMiddleware: authMiddleware}
}

func (oc *OrderController) Route() {
	oc.rg.POST("/ticket/orders", oc.authMiddleware.RequireToken("admin", "Organization"), oc.createOrder)
}

func (oc *OrderController) createOrder(ctx *gin.Context) {
	var payload models.PayloadOrder

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userModel := user.(models.User)
	payload.UserID = userModel.Id

	order, err := oc.orderUseCase.CreateOrder(payload)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success create order ticket", order, true))
}
