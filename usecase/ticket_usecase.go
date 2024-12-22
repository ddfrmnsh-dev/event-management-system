package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
	"fmt"
)

type TicketUseCase interface {
	CreateTicket(input []models.Tickets) ([]models.Tickets, error)
}

type ticketUseCaseImpl struct {
	ticketRepository repository.TicketRepository
}

func NewTicketUseCase(ticketRepository repository.TicketRepository) TicketUseCase {
	return &ticketUseCaseImpl{ticketRepository: ticketRepository}
}

func (tc *ticketUseCaseImpl) CreateTicket(input []models.Tickets) ([]models.Tickets, error) {
	var errs error
	ticket := []models.Tickets{}

	if len(input) == 0 {
		fmt.Println("No tickets to process")
		return []models.Tickets{}, errs
	}

	for _, t := range input {
		newTicket := models.Tickets{
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
