package controller

import (
	"event-management-system/middleware"
	"event-management-system/usecase"
	modelutil "event-management-system/utils/model_util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventUseCase   usecase.EventUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewEventController(eventUseCase usecase.EventUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *EventController {
	return &EventController{eventUseCase: eventUseCase, rg: rg, authMiddleware: authMiddleware}
}

func (ec *EventController) Route() {
	ec.rg.GET("/events", ec.getAllEvent)
	ec.rg.GET("/event/:id", ec.getEventById)
}

func (ec *EventController) getAllEvent(ctx *gin.Context) {
	events, err := ec.eventUseCase.FindAllEvent()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data events"})
		return
	}

	if len(events) > 0 {
		ctx.JSON(http.StatusOK, modelutil.APIResponse("Success get all data event", gin.H{"events": events}, true))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("List event empty", nil, false))
}

func (ec *EventController) getEventById(ctx *gin.Context) {
	eventId := ctx.Param("id")

	id, err := strconv.Atoi(eventId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid event id"})
		return
	}

	event, err := ec.eventUseCase.FindEventById(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success Get Event", event, true))
}
