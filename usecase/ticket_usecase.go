package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
	"fmt"
)

type TicketUseCase interface {
	CreateTicket(input []models.Ticket) ([]models.Ticket, error)
	DeleteTicketById(id []int) (models.Ticket, error)
}

type ticketUseCaseImpl struct {
	ticketRepository repository.TicketRepository
}

func NewTicketUseCase(ticketRepository repository.TicketRepository) TicketUseCase {
	return &ticketUseCaseImpl{ticketRepository: ticketRepository}
}

func (tc *ticketUseCaseImpl) CreateTicket(input []models.Ticket) ([]models.Ticket, error) {
	var errs error
	ticket := []models.Ticket{}

	if len(input) == 0 {
		fmt.Println("No tickets to process")
		return []models.Ticket{}, errs
	}

	for _, t := range input {
		newTicket := models.Ticket{
			TikcetUuid: GenerateUuid(),
			TicketType: t.TicketType,
			Price:      t.Price,
			Quota:      t.Quota,
			Status:     t.Status,
			EventID:    t.EventID,
		}
		ticket = append(ticket, newTicket)
	}

	saveTicket, err := tc.ticketRepository.Save(ticket)
	if err != nil {
		return saveTicket, err
	}

	return saveTicket, nil
}

func (tc *ticketUseCaseImpl) DeleteTicketById(id []int) (models.Ticket, error) {
	ticket, err := tc.ticketRepository.Delete(id)
	if err != nil {
		return ticket, err
	}
	return ticket, nil
}
