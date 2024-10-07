package repository

import (
	"event-management-system/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAllUser(users *[]models.User) error
	FindUserById(id int) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) FindAllUser(users *[]models.User) error {
	res := u.db.Find(&users)
	if res.Error != nil {
		log.Printf(fmt.Sprintf("error fetching all user::%v", res.Error))
		return nil
	}

	return nil
}

func (u *userRepositoryImpl) FindUserById(id int) error {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return err
	}

	return nil
}
