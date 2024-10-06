package repository

import (
	"event-management-system/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUser(users *[]models.User) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u userRepositoryImpl) FindAllUser(users *[]models.User) error {
	res := u.db.Find(&users)
	if res.Error != nil {
		log.Printf(fmt.Sprintf("error fetching all user::%v", res.Error))
		return nil
	}

	return nil
}
