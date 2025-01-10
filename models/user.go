package models

import (
	"time"
)

// User represents the structure of a user.
// @Description User details
type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null;size:255"`
	Email    string `json:"email" gorm:"not null;size:255"`
	Password string `json:"password" gorm:"not null;size:255"`
	// Role      string    `json:"role" gorm:"not null;size:50"`
	IsActive  *bool     `json:"isActive" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"default:current_timestamp"`
	Events    []Event   `gorm:"foreignKey:UserID"`
	Orders    []Order   `gorm:"foreignKey:UserID"`
	Role      []Role    `gorm:"many2many:user_has_roles;constraint:OnDelete:CASCADE"`
}

type PayloadRole struct {
	RoleId       int   `json:"roleId" binding:"required"`
	PermissionId []int `json:"permissionId" binding:"required"`
}

type Role struct {
	ID          int          `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_has_permissions;constraint:OnDelete:CASCADE" json:"permissions"`
	CreatedAt   time.Time    `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time    `json:"updatedAt" gorm:"default:current_timestamp"`
}

type Permission struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"default:current_timestamp"`
}

type GetCustomerDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}

// InputLogin godoc
// @Description Login credentials
// @Property usernameOrEmail string "Identifier for login"
// @Property password string "Password for login"
type InputLogin struct {
	Identifier string `json:"usernameOrEmail" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type UserResponse struct {
	Id        int            `json:"id"`
	Email     string         `json:"email"`
	Username  string         `json:"username"`
	IsActive  *bool          `json:"isActive"`
	Role      []RoleResponse `json:"roles"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// create role response
type RoleResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
}

// create permission response
type PermissionResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// func FormatUserResponse(user User) UserResponse {
// 	return UserResponse{
// 		Id:        user.Id,
// 		Email:     user.Email,
// 		Username:  user.Username,
// 		IsActive:  user.IsActive,
// 		Role:      user.Role,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}
// }

func FormatUserResponse(user User) UserResponse {
	// Konversi roles ke RoleResponse
	var rolesResponse []RoleResponse
	for _, role := range user.Role {
		// Konversi permissions ke PermissionResponse
		var permissionsResponse []PermissionResponse
		for _, permission := range role.Permissions {
			permissionsResponse = append(permissionsResponse, PermissionResponse{
				ID:          permission.ID,
				Name:        permission.Name,
				Description: permission.Description,
			})
		}

		// Tambahkan role ke rolesResponse
		rolesResponse = append(rolesResponse, RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionsResponse,
		})
	}

	// Format user response
	return UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Username:  user.Username,
		IsActive:  user.IsActive,
		Role:      rolesResponse,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
