package controller

import (
	"event-management-system/middleware"
	"event-management-system/models"
	"event-management-system/usecase"
	modelutil "event-management-system/utils/model_util"
	"fmt"
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
	ec.rg.POST("/events", ec.authMiddleware.RequireToken("admin"), ec.createEvent)
	ec.rg.GET("/event/:id", ec.getEventById)
	ec.rg.PUT("/event/:id", ec.authMiddleware.RequireToken("admin"), ec.updateEvent)
	ec.rg.DELETE("/event/:id", ec.authMiddleware.RequireToken("admin"), ec.deleteEvent)
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

func (ec *EventController) createEvent(ctx *gin.Context) {
	var payload models.Event

	if err := ctx.ShouldBind(&payload); err != nil {
		// if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Err JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// file, err := ctx.FormFile("image")
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"err image": err.Error()})
	// 	return
	// }

	// uploadDir := "./uploads"
	// os.MkdirAll(uploadDir, os.ModePerm)
	// fileName := filepath.Base(file.Filename)
	// filePath := filepath.Join(uploadDir, fileName)

	// if err := ctx.SaveUploadedFile(file, filePath); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
	// }

	// payload.PathImage = filePath

	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	userModel := user.(models.User)
	userId := userModel.Id

	payload.UserID = userId
	event, err := ec.eventUseCase.CreateEvent(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	users, err := ec.eventUseCase.FindEventUser(event.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	userResponse := models.FormatUserResponse(users)
	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success create event", gin.H{
		"id":          event.Id,
		"eventUuid":   event.EventUuid,
		"name":        event.Name,
		"slug":        event.Slug,
		"statusEvent": event.StatusEvent,
		"startDate":   event.StartDate,
		"endDate":     event.EndDate,
		"startTime":   event.StartTime,
		"endTime":     event.EndTime,
		"location":    event.Location,
		"address":     event.Address,
		"description": event.Description,
		"ticketTypes": event.TicketTypes,
		"pathImage":   event.PathImage,
		"minPrice":    event.MinimumPrice,
		"createdAt":   event.CreatedAt,
		"updatedAt":   event.UpdatedAt,
		"user":        userResponse,
	}, true))
}

func (ec *EventController) updateEvent(ctx *gin.Context) {
	var inputId models.GetEventDetailInput
	err := ctx.ShouldBindUri(&inputId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input models.Event
	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventUpdated, err := ec.eventUseCase.UpdateEvent(inputId, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success update event", eventUpdated, true))
}

func (ec *EventController) deleteEvent(ctx *gin.Context) {
	var inputId models.GetEventDetailInput
	err := ctx.ShouldBindUri(&inputId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newId, _ := strconv.Atoi(inputId.Id)
	_, err = ec.eventUseCase.DeleteEventById(newId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success delete event", nil, true))
}
