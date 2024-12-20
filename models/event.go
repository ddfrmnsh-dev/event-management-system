package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Event struct {
	Id           int        `json:"id" form:"id" gorm:"primaryKey"`
	EventUuid    string     `json:"eventUuid" form:"eventUuid" gorm:"size:255"`
	Name         string     `json:"name" form:"name" gorm:"not null;size:255"`
	Slug         string     `json:"slug" form:"slug" gorm:"not null;size:255"`
	StatusEvent  string     `json:"statusEvent" form:"statusEvent" gorm:"not null;size:255"`
	StartDate    time.Time  `json:"startDate" form:"startDate" gorm:"type:date;not null"`
	EndDate      time.Time  `json:"endDate" form:"endDate" gorm:"type:date;not null"`
	StartTime    string     `json:"startTime" form:"startTime" gorm:"not null"`
	EndTime      string     `json:"endTime" form:"endTime" gorm:"not null"`
	Location     string     `json:"location" form:"location" gorm:"not null;size:255"`
	Address      string     `json:"address" form:"address" gorm:"not null;size:255"`
	Description  string     `json:"description" form:"description" gorm:"not null"`
	TicketTypes  string     `json:"ticketTypes" form:"ticketTypes" gorm:"not null;size:255"`
	PathImage    string     `json:"pathImage" form:"pathImage" gorm:"not null"`
	MinimumPrice int        `json:"minPrice" form:"minPrice" gorm:"not null"`
	CreatedAt    time.Time  `json:"createdAt" form:"createdAt" gorm:"not null"`
	UpdatedAt    *time.Time `json:"updatedAt" form:"updatedAt" gorm:"autoUpdateTime:false"`
	UserID       int        `json:"userId" gorm:"not null"`
	User         User       `gorm:"constraint:OnDelete:CASCADE;"`
}

// func FormatUserEvent(user User, event Event) Event {
// 	formatter := Event{
// 		Id:           event.Id,
// 		EventUuid:    event.EventUuid,
// 		Name:         event.Name,
// 		Slug:         event.Slug,
// 		StatusEvent:  event.StatusEvent,
// 		StartDate:    event.StartDate,
// 		EndDate:      event.EndDate,
// 		StartTime:    event.StartTime,
// 		EndTime:      event.EndTime,
// 		Location:     event.Location,
// 		Address:      event.Address,
// 		Description:  event.Description,
// 		TicketTypes:  event.TicketTypes,
// 		PathImage:    event.PathImage,
// 		MinimumPrice: event.MinimumPrice,
// 		CreatedAt:    event.CreatedAt,
// 		UpdatedAt:    event.UpdatedAt,
// 		UserID:       event.UserID,
// 		User: User{
// 			Id:        user.Id,
// 			Username:  user.Username,
// 			IsActive:  user.IsActive,
// 			CreatedAt: user.CreatedAt,
// 			UpdatedAt: user.UpdatedAt,
// 		},
// 	}
// 	return formatter
// }

type GetEventDetailInput struct {
	Id string `uri:"id" binding:"required,numeric"`
}

type StatusType string

const (
	Publish  StatusType = "publish"
	Inactive StatusType = "inactive"
	Pending  StatusType = "pending"
	Draft    StatusType = "draft"
)

func (p *StatusType) Scan(value interface{}) error {
	*p = StatusType(value.([]byte))
	return nil
}

func (p StatusType) Value() (driver.Value, error) {
	return string(p), nil
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	const dateFormat = "2006-01-02"

	// loc, err := time.LoadLocation("Asia/Jakarta")
	// if err != nil {
	// 	return fmt.Errorf("failed to load location: %w", err)
	// }

	parsedStartDate, err := time.Parse(dateFormat, aux.StartDate)
	if err != nil {
		return err
	}
	parsedEndDate, err := time.Parse(dateFormat, aux.EndDate)
	if err != nil {
		return err
	}

	// parsedStartTime, err := time.Parse(timeFormat, aux.StartTime)
	// if err != nil {
	// 	return err
	// }
	// parsedEndTime, err := time.Parse(timeFormat, aux.EndTime)
	// if err != nil {
	// 	return err
	// }

	e.StartDate = parsedStartDate
	e.EndDate = parsedEndDate
	// e.StartTime = parsedStartTime
	// e.EndTime = parsedEndTime

	return nil
}
