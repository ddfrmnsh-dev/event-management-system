package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type EventUseCase interface {
	FindAllEvent() ([]models.Event, error)
	FindEventById(id int) (models.Event, error)
	FindEventUser(id int) (models.User, error)
	FindParticipantEvent(id int) (models.Event, error)
	FindEventTicket() ([]models.Event, error)
	// FindEventTicket(id int) (models.Event, error)
	// FindEventTicket(id int) ([]models.Ticket, error)
	CreateEvent(input models.Event) (models.Event, error)
	UpdateEvent(uriId models.GetEventDetailInput, input models.Event) (models.Event, error)
	DeleteEventById(id int) (models.Event, error)
}

type eventUseCaseImpl struct {
	eventRepository repository.EventRepository
	userRepository  repository.UserRepository
}

func NewEventUseCase(eventRepository repository.EventRepository, userRepository repository.UserRepository) EventUseCase {
	return &eventUseCaseImpl{eventRepository: eventRepository, userRepository: userRepository}
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
	event.UserID = input.UserID
	// event.UpdatedAt = nil

	saveEvent, err := ec.eventRepository.Save(event)
	if err != nil {
		return saveEvent, err
	}

	return saveEvent, err
}

func (ec *eventUseCaseImpl) UpdateEvent(uriId models.GetEventDetailInput, input models.Event) (models.Event, error) {
	newId, _ := strconv.Atoi(uriId.Id)
	checkEvent, err := ec.eventRepository.FindById(newId)
	if err != nil {
		return checkEvent, err
	}

	if strings.TrimSpace(input.Name) != "" {
		checkEvent.Name = input.Name
		checkEvent.Slug = slug.Make(input.Name)
	}
	if strings.TrimSpace(input.StatusEvent) != "" {
		checkEvent.StatusEvent = input.StatusEvent
	}
	if strings.TrimSpace(input.StartTime) != "" || strings.TrimSpace(input.EndTime) != "" {
		checkEvent.StartTime = input.StartTime
		checkEvent.EndTime = input.EndTime
	}
	if strings.TrimSpace(input.Location) != "" {
		checkEvent.Location = input.Location
	}
	if strings.TrimSpace(input.Address) != "" {
		checkEvent.Address = input.Address
	}
	if strings.TrimSpace(input.Description) != "" {
		checkEvent.Description = input.Description
	}
	if strings.TrimSpace(input.TicketTypes) != "" {
		checkEvent.TicketTypes = input.TicketTypes
	}
	if strings.TrimSpace(input.PathImage) != "" {
		checkEvent.PathImage = input.PathImage
	}
	if input.MinimumPrice != 0 {
		checkEvent.MinimumPrice = input.MinimumPrice
	}

	timeNow := time.Now()
	checkEvent.UpdatedAt = &timeNow

	updateEvent, err := ec.eventRepository.Update(checkEvent)

	if err != nil {
		return updateEvent, err
	}

	return updateEvent, nil
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

func (ec *eventUseCaseImpl) FindEventUser(id int) (models.User, error) {
	user, err := ec.userRepository.FindById(id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ec *eventUseCaseImpl) DeleteEventById(id int) (models.Event, error) {
	event, err := ec.eventRepository.Delete(id)

	if err != nil {
		return event, err
	}

	return event, nil
}

func (ec *eventUseCaseImpl) FindEventTicket() ([]models.Event, error) {
	event, err := ec.eventRepository.FindEventTicket()

	if err != nil {
		return event, err
	}

	return event, nil
}

func (ec *eventUseCaseImpl) FindParticipantEvent(id int) (models.Event, error) {
	event, err := ec.eventRepository.FindParticipantEvent(id)

	if err != nil {
		return event, err
	}

	return event, nil
}

// func (ec *eventUseCaseImpl) FindEventTicket(id int) ([]models.Ticket, error) {
// 	event, err := ec.eventRepository.FindEventTicket(id)

// 	if err != nil {
// 		return event, err
// 	}

// 	return event, nil
// }

func GenerateUuid() string {
	return uuid.NewString()
}
