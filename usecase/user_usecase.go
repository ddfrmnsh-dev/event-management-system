package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
)

type UserUseCase interface {
	FindAllUser() ([]models.User, error)
	FindUserById(id int) (models.User, error)
}

type userUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCaseImpl{userRepository: userRepository}
}

func (uc *userUseCaseImpl) FindAllUser() ([]models.User, error) {
	var users []models.User
	if err := uc.userRepository.FindAllUser(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *userUseCaseImpl) FindUserById(id int) (models.User, error) {
	var user models.User
	if err := uc.userRepository.FindUserById(id); err != nil {
		return user, err
	}
	return user, nil
}
