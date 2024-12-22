package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Save(ticket []models.Tickets) ([]models.Tickets, error)
}

type ticketRepositoryImpl struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *ticketRepositoryImpl {
	return &ticketRepositoryImpl{db: db}
}

func (t *ticketRepositoryImpl) Save(ticket []models.Tickets) ([]models.Tickets, error) {
	if len(ticket) == 0 {
		return []models.Tickets{}, fmt.Errorf("no tickets to save")
	}
	res := t.db.Create(&ticket)

	if res.Error != nil {
		return ticket, res.Error
	}

	return ticket, nil
}
