package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type EventRepository interface {
	FindAll() ([]models.Event, error)
	FindById(id int) (models.Event, error)
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepositoryImpl {
	return &eventRepositoryImpl{db: db}
}

func (e *eventRepositoryImpl) FindAll() ([]models.Event, error) {
	var events []models.Event

	res := e.db.Find(&events)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("no events found")
		return events, nil
	}

	return events, nil
}

func (e *eventRepositoryImpl) FindById(id int) (models.Event, error) {
	var event models.Event

	res := e.db.First(&event, id)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return event, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("no event found")
		return event, nil
	}

	return event, nil
}
