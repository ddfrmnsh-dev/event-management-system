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
	eventModel     models.Event
}

func NewEventController(eventUseCase usecase.EventUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *EventController {
	return &EventController{eventUseCase: eventUseCase, rg: rg, authMiddleware: authMiddleware}
}

func (ec *EventController) Route() {
	ec.rg.GET("/events", ec.getAllEvent)
	ec.rg.POST("/events", ec.authMiddleware.RequireToken("admin"), ec.createEvent)
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
	event, err := ec.eventUseCase.CreateEvent(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	fmt.Println("data", event)
	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success create event", event, true))
}

// func (ec *EventController) createEvent(ctx *gin.Context) {
// 	var payload models.Event

// 	// Bind data dari form-data
// 	if err := ctx.ShouldBind(&payload); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %s", err.Error())})
// 		return
// 	}

// 	// Parsing tanggal dan waktu dari form-data
// 	const dateFormat = "2006-01-02"
// 	const dateTimeFormat = "2006-01-02T15:04:05Z07:00"
// 	loc, _ := time.LoadLocation("Asia/Jakarta")

// 	// Parse StartDate
// 	if startDateStr := ctx.PostForm("startDate"); startDateStr != "" {
// 		startDate, err := time.ParseInLocation(dateTimeFormat, startDateStr, loc)
// 		if err != nil {
// 			// Jika format tidak cocok dengan format lengkap, coba dengan format tanggal saja
// 			startDate, err = time.ParseInLocation(dateFormat, startDateStr, loc)
// 			if err != nil {
// 				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid startDate: %s", err.Error())})
// 				return
// 			}
// 		}
// 		payload.StartDate = startDate
// 	}

// 	// Parse EndDate
// 	if endDateStr := ctx.PostForm("endDate"); endDateStr != "" {
// 		endDate, err := time.ParseInLocation(dateTimeFormat, endDateStr, loc)
// 		if err != nil {
// 			// Jika format tidak cocok dengan format lengkap, coba dengan format tanggal saja
// 			endDate, err = time.ParseInLocation(dateFormat, endDateStr, loc)
// 			if err != nil {
// 				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid endDate: %s", err.Error())})
// 				return
// 			}
// 		}
// 		payload.EndDate = endDate
// 	}

// 	// Parse StartTime
// 	if startTimeStr := ctx.PostForm("startTime"); startTimeStr != "" {
// 		startTime, err := time.ParseInLocation("15:04", startTimeStr, loc)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid startTime: %s", err.Error())})
// 			return
// 		}
// 		payload.StartTime = startTime
// 	}

// 	// Parse EndTime
// 	if endTimeStr := ctx.PostForm("endTime"); endTimeStr != "" {
// 		endTime, err := time.ParseInLocation("15:04", endTimeStr, loc)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid endTime: %s", err.Error())})
// 			return
// 		}
// 		payload.EndTime = endTime
// 	}

// 	// Handle file upload (image)
// 	file, err := ctx.FormFile("image")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
// 		return
// 	}

// 	const uploadDir = "./uploads"
// 	os.MkdirAll(uploadDir, os.ModePerm)
// 	fileName := filepath.Base(file.Filename)
// 	filePath := filepath.Join(uploadDir, fileName)

// 	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded image"})
// 		return
// 	}

// 	payload.PathImage = filePath

// 	// Panggil use case untuk membuat event
// 	event, err := ec.eventUseCase.CreateEvent(payload)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create event: %s", err.Error())})
// 		return
// 	}

// 	// Return response berhasil
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Event successfully created",
// 		"data":    event,
// 	})
// }
