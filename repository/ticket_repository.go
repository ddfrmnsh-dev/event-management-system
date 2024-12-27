package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type TicketRepository interface {
	Save(ticket []models.Ticket) ([]models.Ticket, error)
	FindById(id int) ([]models.Ticket, error)
	Delete(idTicket []int) (models.Ticket, error)
}

type ticketRepositoryImpl struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *ticketRepositoryImpl {
	return &ticketRepositoryImpl{db: db}
}

func (t *ticketRepositoryImpl) Save(ticket []models.Ticket) ([]models.Ticket, error) {
	if len(ticket) == 0 {
		return []models.Ticket{}, fmt.Errorf("no tickets to save")
	}
	res := t.db.Create(&ticket)

	if res.Error != nil {
		return ticket, res.Error
	}

	return ticket, nil
}

func (t *ticketRepositoryImpl) Delete(idTicket []int) (models.Ticket, error) {
	var ticket models.Ticket
	if len(idTicket) == 0 {
		return ticket, fmt.Errorf("no IDs provided")
	}

	res := t.db.Where("id IN ?", idTicket).Delete(&ticket)
	if res.Error != nil {
		return ticket, fmt.Errorf("failed to delete items: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return ticket, fmt.Errorf("no items were deleted")
	}

	return ticket, nil
}

func (t *ticketRepositoryImpl) FindById(id int) ([]models.Ticket, error) {
	var ticket []models.Ticket

	res := t.db.Where("event_id = ?", id).Find(&ticket)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return ticket, res.Error
	}

	if res.RowsAffected == 0 {
		return ticket, nil
	}

	return ticket, nil
}
