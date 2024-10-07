package models

import "time"

type User struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null;size:255"`
	Email     string    `json:"email" gorm:"not null;size:255"`
	Password  string    `json:"password" gorm:"not null;size:255"`
	Role      string    `json:"role" gorm:"not null;size:50"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"default:current_timestamp"`
}
