package models

import "time"

type Tickets struct {
	Id         int        `json:"id" form:"id" gorm:"primaryKey"`
	TikcetUuid string     `json:"ticketUuid" form:"ticketUuid" gorm:"size:255"`
	TicketType string     `json:"ticketType"`
	Price      int        `json:"price"`
	Quota      int        `json:"quota"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt" form:"updatedAt" gorm:"autoUpdateTime:false"`
	EventID    int        `json:"eventId" gorm:"not null"`
	Event      Event      `gorm:"constraint:OnDelete:CASCADE;"`
}
