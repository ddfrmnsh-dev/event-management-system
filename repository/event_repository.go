package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type EventRepository interface {
	FindAll() ([]models.Event, error)
	FindById(id int) (models.Event, error)
	// FindEventTicket(id int) (models.Event, error)
	FindEventTicket() ([]models.Event, error)
	FindParticipantEvent(id int) (models.Event, error)
	// FindEventTicket(id int) ([]models.Ticket, error)
	Save(event models.Event) (models.Event, error)
	Update(event models.Event) (models.Event, error)
	Delete(id int) (models.Event, error)
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

	res := e.db.Preload("Tickets").First(&event, id)

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

func (e *eventRepositoryImpl) Save(event models.Event) (models.Event, error) {
	res := e.db.Create(&event)
	fmt.Println("res", event)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return event, res.Error
	}

	// if err := e.db.Preload("User").First(&event, event.Id).Error; err != nil {
	// 	fmt.Println("err:", res.Error)
	// 	return event, res.Error
	// }

	return event, nil
}

func (e *eventRepositoryImpl) Update(event models.Event) (models.Event, error) {
	res := e.db.Updates(&event)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return event, res.Error
	}

	return event, nil
}

func (e *eventRepositoryImpl) Delete(id int) (models.Event, error) {
	checkId, err := e.FindById(id)

	if err != nil {
		return checkId, err
	}

	var event models.Event
	res := e.db.Delete(&event, id)

	if res.Error != nil {
		return event, res.Error
	}

	return checkId, nil
}

// func (e *eventRepositoryImpl) FindEventTicket(id int) ([]models.Ticket, error) {
// 	// var event models.Event
// 	// res := e.db.Preload("Tickets").Where("id = ?", id).First(&event)
// 	// fmt.Printf("Query executed: %+v\n", res.Statement.SQL.String())

// 	// if res.Error != nil {
// 	// 	fmt.Println("err:", res.Error)
// 	// 	return event, res.Error
// 	// }

// 	// if res.RowsAffected == 0 {
// 	// 	fmt.Println("no event found")
// 	// 	return event, nil
// 	// }

// 	// return event, nil

// 	var tickets []models.Ticket
// 	res := e.db.Where("event_id = ?", id).Find(&tickets)
// 	if res.Error != nil {
// 		return nil, res.Error
// 	}
// 	return tickets, nil

// }

func (e *eventRepositoryImpl) FindEventTicket() ([]models.Event, error) {
	var event []models.Event
	res := e.db.Preload("Tickets").Find(&event)
	fmt.Printf("Query executed: %+v\n", res.Statement.SQL.String())

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
func (e *eventRepositoryImpl) FindParticipantEvent(id int) (models.Event, error) {
	var event models.Event
	res := e.db.Preload("User.Orders.OrderDetails.Ticket").First(&event, id)
	fmt.Printf("Query executed: %+v\n", res.Statement.SQL.String())

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

// func (e *eventRepositoryImpl) FindEventTicket(id int) (models.Event, error) {
// 	var event models.Event
// 	res := e.db.Preload("Tickets").Where("id = ?", id).Find(&event)
// 	fmt.Printf("Query executed: %+v\n", res.Statement.SQL.String())

// 	if res.Error != nil {
// 		fmt.Println("err:", res.Error)
// 		return event, res.Error
// 	}

// 	if res.RowsAffected == 0 {
// 		fmt.Println("no event found")
// 		return event, nil
// 	}

// 	return event, nil

// }
