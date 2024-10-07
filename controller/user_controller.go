package controller

import (
	"event-management-system/usecase"
	"net/http"
	"strconv"

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
	uc.rg.GET("/users/:id", uc.getUserById)
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

func (uc *UserController) getUserById(ctx *gin.Context) {
	idUser := ctx.Param("id")

	id, err := strconv.Atoi(idUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid id"})
		return
	}
	user, err := uc.userUseCase.FindUserById(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
