package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
)

type EventUseCase interface {
	FindAllEvent() ([]models.Event, error)
	FindEventById(id int) (models.Event, error)
}

type eventUseCaseImpl struct {
	eventRepository repository.EventRepository
}

func NewEventUseCase(eventRepository repository.EventRepository) EventUseCase {
	return &eventUseCaseImpl{eventRepository: eventRepository}
}

func (ec *eventUseCaseImpl) FindAllEvent() ([]models.Event, error) {
	events, err := ec.eventRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (ec *eventUseCaseImpl) FindEventById(id int) (models.Event, error) {
	event, err := ec.eventRepository.FindById(id)

	if err != nil {
		return event, err
	}

	return event, nil
}
