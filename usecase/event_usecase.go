package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type EventUseCase interface {
	FindAllEvent() ([]models.Event, error)
	FindEventById(id int) (models.Event, error)
	CreateEvent(input models.Event) (models.Event, error)
}

type eventUseCaseImpl struct {
	eventRepository repository.EventRepository
}

func NewEventUseCase(eventRepository repository.EventRepository) EventUseCase {
	return &eventUseCaseImpl{eventRepository: eventRepository}
}

func (ec *eventUseCaseImpl) CreateEvent(input models.Event) (models.Event, error) {
	event := models.Event{}

	event.EventUuid = GenerateUuid()
	event.Name = input.Name
	event.Slug = slug.Make(input.Name)
	event.StatusEvent = input.StatusEvent
	event.StartDate = input.StartDate
	event.EndDate = input.EndDate
	event.StartTime = input.StartTime
	event.EndTime = input.EndTime
	event.Location = input.Location
	event.Address = input.Address
	event.Description = input.Description
	event.TicketTypes = input.TicketTypes
	event.PathImage = input.PathImage
	event.MinimumPrice = input.MinimumPrice
	// event.UpdatedAt = nil

	saveEvent, err := ec.eventRepository.Save(event)
	if err != nil {
		return saveEvent, err
	}

	return saveEvent, err
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

func GenerateUuid() string {
	return uuid.NewString()
}
