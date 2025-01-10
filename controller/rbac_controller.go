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

type RbacController struct {
	rbacUseCase    usecase.RbacUseCase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewRbacController(rbacUseCase usecase.RbacUseCase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *RbacController {
	return &RbacController{rbacUseCase: rbacUseCase, rg: rg, authMiddleware: authMiddleware}
}

func (rc *RbacController) Route() {
	rc.rg.GET("/roles", rc.authMiddleware.RequireToken("admin", "super admin"), rc.getAllRole)
	rc.rg.GET("/permissions", rc.authMiddleware.RequireToken("admin", "super admin"), rc.getAllPermission)
	rc.rg.POST("/roles", rc.authMiddleware.RequireToken("admin", "super admin"), rc.createRole)
	rc.rg.POST("/permissions", rc.authMiddleware.RequireToken("admin", "super admin"), rc.createPermission)
	rc.rg.DELETE("/role/:id", rc.authMiddleware.RequireToken("admin", "super admin"), rc.deleteRole)
	rc.rg.DELETE("/permission/:id", rc.authMiddleware.RequireToken("admin", "super admin"), rc.deletePermission)
	rc.rg.POST("/role/assignPermission", rc.authMiddleware.RequireToken("admin", "super admin"), rc.addPermissionToRole)
}

func (rc *RbacController) getAllRole(ctx *gin.Context) {
	roles, err := rc.rbacUseCase.FindAllRole()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data roles"})
		return
	}

	if len(roles) > 0 {
		ctx.JSON(http.StatusOK, modelutil.APIResponse("Success get all data role", gin.H{"roles": roles}, true))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("List role empty", nil, false))
}
func (rc *RbacController) getAllPermission(ctx *gin.Context) {
	permissions, err := rc.rbacUseCase.FindAllPermission()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to retrieve data permissions"})
		return
	}

	if len(permissions) > 0 {
		ctx.JSON(http.StatusOK, modelutil.APIResponse("Success get all data role", gin.H{"permissions": permissions}, true))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("List role empty", nil, false))
}

func (rc *RbacController) createRole(ctx *gin.Context) {
	var payload models.Role

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Err JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	role, err := rc.rbacUseCase.CreateRole(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success create role", role, true))
}

func (rc *RbacController) createPermission(ctx *gin.Context) {
	var payload models.Permission

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Err JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	permission, err := rc.rbacUseCase.CreatePermission(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success create permission", permission, true))
}

func (rc *RbacController) deleteRole(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, modelutil.APIResponse("Invalid ID", nil, false))
		return
	}

	_, err = rc.rbacUseCase.DeleteRole(idInt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}
	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success delete role", nil, true))
}
func (rc *RbacController) deletePermission(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, modelutil.APIResponse("Invalid ID", nil, false))
		return
	}

	_, err = rc.rbacUseCase.DeletePermission(idInt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}
	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success delete permission", nil, true))
}

func (rc *RbacController) addPermissionToRole(ctx *gin.Context) {
	var payload models.PayloadRole

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Err JSON:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	role, err := rc.rbacUseCase.AddPermissionToRole(payload.RoleId, payload.PermissionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, modelutil.APIResponse(err.Error(), nil, false))
		return
	}

	ctx.JSON(http.StatusOK, modelutil.APIResponse("Success add permission to role", role, true))
}
