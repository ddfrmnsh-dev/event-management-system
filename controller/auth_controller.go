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
	ac.rg.POST("/signinAdmin", ac.loginAdmin)
	ac.rg.POST("/signinUser", ac.loginUser)
}

// Login godoc
// @Summary Login Admin User
// @Description Authenticate user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.InputLogin true "Login credentials"
// @Success 200 {object} modelUtil.Response
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/signinAdmin [post]
func (ac *AuthController) loginAdmin(ctx *gin.Context) {
	fmt.Println("Starting login process")
	var payload models.InputLogin

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("Cek err", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input"})
		return
	}

	token, user, err := ac.authUseCase.LoginAdmin(payload.Identifier, payload.Password)

	if err != nil {
		if err.Error() == "invalid password" || err.Error() == "user not found" || err.Error() == "invalid credentials" {
			ctx.JSON(http.StatusUnauthorized, modelUtil.APIResponse("Invalid credentials", nil, false))
			return
		}

		fmt.Println("Cek err", err.Error())
		ctx.JSON(http.StatusInternalServerError, modelUtil.APIResponse("Failed to process request", nil, false))
		return
	}

	formatedUser := models.FormatUserResponse(user)
	ctx.JSON(http.StatusOK, modelUtil.APIResponse("success", gin.H{
		"userPrincipal": formatedUser,
		"accessToken":   token,
	}, true))
}

// Login godoc
// @Summary Login User
// @Description Authenticate user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.InputLogin true "Login credentials"
// @Success 200 {object} modelUtil.Response
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/signinUser [post]
func (ac *AuthController) loginUser(ctx *gin.Context) {
	fmt.Println("Starting login process")
	var payload models.InputLogin

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		fmt.Println("Cek err", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input"})
		return
	}

	token, user, err := ac.authUseCase.LoginUser(payload.Identifier, payload.Password)

	if err != nil {
		if err.Error() == "invalid password" || err.Error() == "user not found" || err.Error() == "invalid credentials" {
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
