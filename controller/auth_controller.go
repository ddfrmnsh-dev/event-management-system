package controller

import (
	"event-management-system/models"
	"event-management-system/usecase"
	modelUtil "event-management-system/utils/model_util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase usecase.AuthenticationUseCase
	rg          *gin.RouterGroup
}

func NewAuthController(authUseCase usecase.AuthenticationUseCase, rg *gin.RouterGroup) *AuthController {
	return &AuthController{authUseCase: authUseCase, rg: rg}
}

func (ac *AuthController) Route() {
	ac.rg.POST("/signin", ac.login)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.InputLogin true "Login credentials"
// @Success 200 {object} modelUtil.Response
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/signin [post]
func (ac *AuthController) login(ctx *gin.Context) {
	fmt.Println("Starting login process")
	var payload models.InputLogin

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("Cek err", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input"})
		return
	}

	token, user, err := ac.authUseCase.Login(payload.Identifier, payload.Password)

	if err != nil {
		if err.Error() == "invalid password" || err.Error() == "user not found" {
			ctx.JSON(http.StatusUnauthorized, modelUtil.APIResponse("Invalid credentials", nil, false))
			return
		}

		fmt.Println("Cek err", err.Error())
		ctx.JSON(http.StatusInternalServerError, modelUtil.APIResponse("Failed to process request", nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelUtil.APIResponse("success", gin.H{
		"user":        user,
		"accessToken": token,
	}, true))
}
