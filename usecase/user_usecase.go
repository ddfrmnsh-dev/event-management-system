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
	users, err := uc.userRepository.FindAll() // Mengambil semua pengguna dari repository
	if err != nil {
		return nil, err // Mengembalikan error jika ada
	}

	return users, nil // Mengembalikan slice pengguna
}

func (uc *userUseCaseImpl) FindUserById(id int) (models.User, error) {
	user, err := uc.userRepository.FindById(id)
	if err != nil {
		return user, err
	}

	return user, nil
}
