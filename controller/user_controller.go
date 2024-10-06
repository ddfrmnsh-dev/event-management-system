package controller

import (
	"event-management-system/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase usecase.UserUseCase
	rg          *gin.RouterGroup
}

func NewUserController(userUseCase usecase.UserUseCase, rg *gin.RouterGroup) *UserController {
	return &UserController{userUseCase: userUseCase, rg: rg}
}

func (uc *UserController) Route() {
	uc.rg.GET("/users", uc.getAllUser)
}

func (uc *UserController) getAllUser(ctx *gin.Context) {
	users, err := uc.userUseCase.FindAllUser()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data books"})
		return
	}

	if len(users) > 0 {
		ctx.JSON(http.StatusOK, users)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "List user empty"})
}
