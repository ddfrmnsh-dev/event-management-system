package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type OrderRepository interface {
	BeginTransaction() *gorm.DB
	Save(order models.Order) (models.Order, error)
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepositoryImpl {
	return &orderRepositoryImpl{db: db}
}

func (o *orderRepositoryImpl) BeginTransaction() *gorm.DB {
	return o.db.Begin()
}

func (o *orderRepositoryImpl) Save(order models.Order) (models.Order, error) {
	res := o.db.Create(&order)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return order, res.Error
	}

	return order, nil
}
