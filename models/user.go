package models

import (
	"time"
)

type User struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null;size:255"`
	Email     string    `json:"email" gorm:"not null;size:255"`
	Password  string    `json:"password" gorm:"not null;size:255"`
	Role      string    `json:"role" gorm:"not null;size:50"`
	IsActive  *bool     `json:"isActive" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"default:current_timestamp"`
}

type GetCustomerDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}

type InputLogin struct {
	Identifier string `json:"usernameOrEmail" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
