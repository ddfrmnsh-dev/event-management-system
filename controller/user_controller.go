package controller

import (
	"event-management-system/middleware"
	"event-management-system/models"
	"event-management-system/usecase"
	modelUtil "event-management-system/utils/model_util"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase    usecase.UserUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewUserController(userUseCase usecase.UserUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *UserController {
	return &UserController{userUseCase: userUseCase, rg: rg, authMiddleware: authMiddleware}
}

func (uc *UserController) Route() {
	uc.rg.GET("/users", uc.authMiddleware.RequireToken("admin"), uc.getAllUser)
	uc.rg.GET("/users/:id", uc.authMiddleware.RequireToken("admin"), uc.getUserById)
	uc.rg.POST("/users", uc.authMiddleware.RequireToken("admin"), uc.createUser)
	uc.rg.PUT("/users/:id", uc.authMiddleware.RequireToken("admin"), uc.updateUser)
	uc.rg.DELETE("/users/:id", uc.authMiddleware.RequireToken("admin"), uc.deleteUser)
}

func (uc *UserController) getAllUser(ctx *gin.Context) {
	users, err := uc.userUseCase.FindAllUser()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data users"})
		return
	}

	if len(users) > 0 {
		ctx.JSON(http.StatusOK, modelUtil.APIResponse("Success get all data user", gin.H{"users": users}, true))
		return
	}

	ctx.JSON(http.StatusOK, modelUtil.APIResponse("List user empty", nil, false))
}

func (uc *UserController) getUserById(ctx *gin.Context) {
	idUser := ctx.Param("id")

	id, err := strconv.Atoi(idUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid id"})
		return
	}
	user, err := uc.userUseCase.FindUserById(id)

	fmt.Println(user)
	if err != nil {
		ctx.JSON(http.StatusNotFound, modelUtil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelUtil.APIResponse("Succes Get User", user, true))
}

func (uc *UserController) createUser(ctx *gin.Context) {
	var payload models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("Cek err", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user, err := uc.userUseCase.CreateUser(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelUtil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelUtil.APIResponse("Succes create user", user, true))
}

func (uc *UserController) updateUser(ctx *gin.Context) {
	var inputId models.GetCustomerDetailInput
	err := ctx.ShouldBindUri(&inputId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input models.User
	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateCustomer, err := uc.userUseCase.UpdateUser(inputId, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, modelUtil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelUtil.APIResponse("Succes update user", updateCustomer, true))

}

func (uc *UserController) deleteUser(ctx *gin.Context) {
	var inputId models.GetCustomerDetailInput
	err := ctx.ShouldBindUri(&inputId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newId, _ := strconv.Atoi(inputId.Id)
	deleteCustomer, err := uc.userUseCase.DeleteUserById(newId)
	if err != nil {
		log.Println("Terjadi kesalahan:", err)
		ctx.JSON(http.StatusBadRequest, modelUtil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelUtil.APIResponse("Succes delete user "+strconv.Itoa(deleteCustomer.Id), nil, true))
}
