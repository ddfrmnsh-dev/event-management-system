package models

import "time"

type Event struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Slug        string    `json:"slug" gorm:"not null;size:255"`
	Date        time.Time `json:"date" gorm:"not null"`
	Location    string    `json:"location" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	TicketTypes string    `json:"ticketTypes" gorm:"not null"`
}

type GetEventDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}
