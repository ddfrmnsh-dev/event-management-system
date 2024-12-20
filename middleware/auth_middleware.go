package middleware

import (
	"event-management-system/models"
	"event-management-system/utils/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type authHedaer struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var aH authHedaer
		err := ctx.ShouldBindHeader(&aH)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauhtorized"})
			return
		}

		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)

		tokenClaim, err := a.jwtService.VerifyToken(token)

		newId, _ := strconv.Atoi(tokenClaim.UserId)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauhtorized"})
			return
		}

		ctx.Set("user", models.User{Id: newId, Role: tokenClaim.Role})
		validRole := false

		for _, role := range roles {
			if role == tokenClaim.Role {
				validRole = true
				break
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource"})
			return
		}

		ctx.Next()
	}
}
