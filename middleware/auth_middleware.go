package middleware

import (
	"event-management-system/models"
	"event-management-system/utils/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthMiddleware interface {
	RequireToken(allowedRoles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
	db         *gorm.DB
}

type authHedaer struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

func NewAuthMiddleware(jwtService service.JwtService, db *gorm.DB) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService, db: db}
}

func (a *authMiddleware) RequireToken(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind header
		var aH authHedaer
		if err := ctx.ShouldBindHeader(&aH); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "status": false})
			return
		}

		// Extract token
		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)

		// Verify token
		tokenClaim, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "status": false})
			return
		}

		// Parse UserID
		userId, err := strconv.Atoi(tokenClaim.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid user ID in token", "status": false})
			return
		}

		// Fetch user with roles
		var user models.User
		if err := a.db.Preload("Role").First(&user, userId).Error; err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not found", "status": false})
			return
		}

		// Check if user has any allowed role
		validRole := false
		for _, userRole := range user.Role {
			for _, allowedRole := range allowedRoles {
				if strings.ToLower(userRole.Name) == allowedRole {
					validRole = true
					break
				}
			}
			if validRole {
				break
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource", "status": false})
			return
		}

		// Set user in context for downstream use
		ctx.Set("user", user)
		ctx.Next()
	}
}

// func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var aH authHedaer
// 		err := ctx.ShouldBindHeader(&aH)

// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauhtorized", "status": false})
// 			return
// 		}

// 		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)

// 		tokenClaim, err := a.jwtService.VerifyToken(token)

// 		newId, _ := strconv.Atoi(tokenClaim.UserId)

// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauhtorized", "status": false})
// 			return
// 		}

// 		ctx.Set("user", models.User{Id: newId, Role: tokenClaim.Role})
// 		validRole := false

// 		for _, role := range roles {
// 			if role == tokenClaim.Role {
// 				validRole = true
// 				break
// 			}
// 		}

// 		if !validRole {
// 			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource", "status": false})
// 			return
// 		}

// 		ctx.Next()
// 	}
// }

// func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		var aH authHedaer
// 		err := ctx.ShouldBindHeader(&aH)

// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauhtorized", "status": false})
// 			return
// 		}

// 		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)

// 		tokenClaim, err := a.jwtService.VerifyToken(token)

// 		newId, _ := strconv.Atoi(tokenClaim.UserId)

// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauhtorized", "status": false})
// 			return
// 		}

// 		ctx.Set("user", models.User{Id: newId, Role: tokenClaim.Role})
// 		validRole := false

// 		for _, role := range roles {
// 			if role == tokenClaim.Role {
// 				validRole = true
// 				break
// 			}
// 		}

// 		if !validRole {
// 			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource", "status": false})
// 			return
// 		}

// 		ctx.Next()
// 	}
// }
